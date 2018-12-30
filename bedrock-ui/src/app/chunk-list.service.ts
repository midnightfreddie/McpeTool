import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ChunkListService {

  constructor(private http:HttpClient) {
    // this.nonChunkKeys = [
    //   'mVillages',
    //   'AutonomousEntities',
    //   'BiomeData',
    //   'Overworld',
    //   'dimension0',
    //   'portals'
    // ];
  }
  keyList = [];
  otherKeys;
  nonChunkKeys = [
    'mVillages',
    'AutonomousEntities',
    'BiomeData',
    'Overworld',
    'dimension0',
    'portals'
  ];
  players = {};
  chunks: { [id: string]: string } = {};
  knownKeys = {};
  unKnownKeys = {};
  entityChunkList = [];

  // Have to use arror function notation to preserve 'this' in success function
  RefreshKeys = () => {
    this.http.get('http://127.0.0.1:8080/api/v1/db/')
      .subscribe(this.RefreshKeysSuccess);
  }

  RefreshKeysSuccess = (data) => {
    let chunk;
    this.keyList = data.keys;
    this.players = {};
    this.chunks = {};
    this.knownKeys = {};
    this.unKnownKeys = {};
    this.entityChunkList = [];
    this.keyList.forEach(e => {
      // regex to match player string keys TODO: use hex key match instead?
      if (/^(~local_player|^player_)/.test(e.stringKey)) {
        this.players[e.stringKey] = e;
      } else if (this.nonChunkKeys.includes(e.stringKey)) {
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

  HttpFail(response) {
    console.error('http failure', response.status, response.statustext, response.data);
  }
}
