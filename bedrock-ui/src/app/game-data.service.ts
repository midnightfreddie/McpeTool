import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

const apiRoot = 'http://127.0.0.1:8080/api/v1';
// Matches strings that are hex strings of all-printable ascii (plus 0x7f, but oh well)
const allPrintableAscii = /^([2-7][0-9a-f])+$/i;
// Matches hex strings for ~local_player and those begining with player_
const playerKeyRegEx = /^(7e6c6f63616c5f706c61796572$|706c617965725f)/i;
const chunkKeyRegEx = /^(?<cx>[\da-f]{8})(?<cz>[\da-f]{8})(?<dimension>[\da-f]{8})?(?<tag>[\da-f]{2})(?<subchunk>[\da-f]{2})?$/i;
const chunkTags = {
  45: 'Data2D',
  46: 'Data2DLegacy',
  47: 'SubChunkPrefix',
  48: 'LegacyTerrain',
  49: 'BlockEntity',
  50: 'Entity',
  51: 'PendingTicks',
  52: 'BlockExtraData',
  53: 'BiomeState',
  54: 'FinalizedState',
  118: 'Version',
}

export const xPerChunk = 16;
export const zPerChunk = 16;
export const yPerPEChunk = 128;
export const yPerSubChunk = 16;

@Injectable({
  providedIn: 'root'
})
export class GameDataService {

  constructor(private http:HttpClient) { }

  keyList: string[] = [];

  // Fetches key list from API and populates class property
  // TODO: Handle errors
  refreshKeys = () => {
    this.http.get(`${apiRoot}/db/`)
      .subscribe( data => this.keyList = data['keys'].map(e => e.hexKey));
  }

  // Returns ASCII strings of hex keys that are all printable values
  stringKeys = () => this.keyList.filter(e => allPrintableAscii.test(e)).map(e => this.asciiHexStringToString(e));

  // Given hex string, returns printable ASCII string; does not validate input
  private asciiHexStringToString = hexString => {
    let out: string = "";
    for (let i=0; i < hexString.length; i +=2) {
      out += String.fromCharCode(parseInt(hexString.substring(i, i+2), 16));
    }
    return out;
  }

  playerKeys = () => this.keyList.filter(e => playerKeyRegEx.test(e)).map(e => this.asciiHexStringToString(e));

  // Returns keys that the author knows of but hasn't handled yet
  knownKeys = () => this.keyList.filter(e => [
    'mVillages',
    'AutonomousEntities',
    'BiomeData',
    'Overworld',
    'dimension0',
    'portals'
  ].includes(this.asciiHexStringToString(e))).map(e => this.asciiHexStringToString(e));

  // Returns keys for chunks. The chunkKeyRegEx matches some of the string keys, so printable ASCII keys
  chunkKeys = () => this.keyList.filter(e => chunkKeyRegEx.test(e) && !allPrintableAscii.test(e));

  // Given a chunk hex key, returns derivable info
  // When passing the output of overworldChunks(), have to add a fake tag for now, e.g. '32'
  chunkKeyInfo = (hexKey: string) => {
    const chunk = hexKey.match(chunkKeyRegEx);
    let out = {};
    out['dimension'] = (chunk['groups'] || {})['dimension'] ? parseInt(chunk['groups']['dimension'], 16) : 0;
    for (let e of ['cx', 'cz', 'tag', 'subchunk']) {
      if ((chunk['groups'] || {})[e]) {
        out[e] = this.hexStringToInt(chunk['groups'][e]);
      }
    }
    return out;
  }

  // Return deduplicated list of 16-length hex strings for overworld chunks
  overworldChunks = () => this.chunkKeys()
    .filter(e => {
      const chunkInfo = this.chunkKeyInfo(e);
      return chunkInfo['dimension'] == 0;
    })
    .map(e => e.substring(0,16))
    .filter((e, i, arr) => arr.indexOf(e) === i);

  private hexStringToInt = (hexString: string): number  => {
    let msbString: string = "";
    for (let i=hexString.length - 2; i >= 0 ; i -=2) {
      msbString += hexString.substring(i, i+2);
    }
    let out = parseInt(msbString, 16);
    // Hard coding for signed 32-bit integers; effectively shorter ints unsigned
    return out > 0x80000000 ? out - 0x100000000 : out;

  }
}
