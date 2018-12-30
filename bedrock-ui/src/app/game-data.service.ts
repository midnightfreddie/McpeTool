import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

const apiRoot = 'http://127.0.0.1:8080/api/v1';
/*
const nonChunkKeys = [
  'mVillages',
  'AutonomousEntities',
  'BiomeData',
  'Overworld',
  'dimension0',
  'portals'
];
*/

@Injectable({
  providedIn: 'root'
})
export class GameDataService {

  constructor(private http:HttpClient) { }

  keyList: string[];
/*
  otherKeys;
  players = {};
  chunks = {};
  knownKeys = {};
  unKnownKeys = {};
  entityChunkList = [];
*/
/*
  zeroKeys = () => {
    this.keyList = [];
    this.players = {};
    this.chunks = {};
    this.knownKeys = {};
    this.unKnownKeys = {};
    this.entityChunkList = [];
  }
*/

  // TODO: Handle errors
  refreshKeys = () => {
    this.http.get(`${apiRoot}/db/`)
      .subscribe( data => this.keyList = data.keys.map(e => e.hexKey));
  }
    /* *** Will move this logic to other methods ***
  // Have to use arrow function notation to preserve 'this' in success function
  refreshKeysSuccess = (data) => {
    this.zeroKeys();
    this.keyList = data.keys.map(e => e.hexKey);
    console.log(this.keyList);
    this.keyList.forEach(e => {
      let chunk: string;
      // regex to match player string keys TODO: use hex key match instead?
      if (/^(~local_player|^player_)/.test(e.stringKey)) {
        this.players[e.stringKey] = e;
      } else if (nonChunkKeys.includes(e.stringKey)) {
        this.knownKeys[e.stringKey] = e;
      // matches 9, 10, 13 or 14 hex digits and captures matches
      // 1: x, 2: z, 3: (if present) dimension, 4: tag / data type, 5: subchunk
      } else if (chunk = e.hexKey.match(/^([\da-f]{8})([\da-f]{8})([\da-f]{8})?([\da-f]{2})([\da-f]{2})?$/i)) {
        // Currently using x, z, and dimension (if present) portion of key to identify the chunk
        let hexKey = chunk[1] + chunk[2] + (chunk[3] || "");
        // Create key if not present
        if (! this.chunks[hexKey]) {
          // create array so buffer.slice is available
          let byteArrayKey = new Uint8Array(hexKey.length / 2);
          // populate the array
          for (let i = 0; i < hexKey.length; i += 2) {
            byteArrayKey[i/2] = parseInt("0x" + hexKey.substr(i, 2), 16);
          }
          this.chunks[hexKey] = {
            dimension: chunk[3] ? new Int32Array((new Uint8Array(byteArrayKey)).buffer.slice(8,12))[0] : 0,
            "x": new Int32Array((new Uint8Array(byteArrayKey)).buffer.slice(0,4))[0],
            "z": new Int32Array((new Uint8Array(byteArrayKey)).buffer.slice(4,8))[0],
            subChunks: [],
            unknown: []
          }
        }
        switch (chunk[4]) {
          case '2d':
            this.chunks[hexKey]['data2d'] = e.url;
            break;
          case '2f':
            this.chunks[hexKey].subChunks.push(e.hexKey);
            break;
            case '31':
            this.chunks[hexKey]['block-entities'] = e.hexKey;
            break;
            case '32':
            this.chunks[hexKey]['entities'] = e.hexKey;
            this.entityChunkList.push(e.hexKey);
            break;
          case '36':
            this.chunks[hexKey]['0x36'] = e.hexKey;
            break;
          case '76':
            this.chunks[hexKey]['version'] = e.hexKey;
            break;
          default:
            this.chunks[hexKey].unknown.push(e.hexKey);
        }
      } else {
        this.unKnownKeys[e.hexKey] = e.hexKey;
      }
    });
  }
    */
/*
  httpFail = (response) => {
    console.error('http failure', response.status, response.statustext, response.data);
  }
*/
}
