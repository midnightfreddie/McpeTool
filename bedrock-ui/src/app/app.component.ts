import { Component, OnInit } from '@angular/core';

import { GameDataService } from './game-data.service';

@Component({
  selector: 'app-root',
  template: `
    Under Construction {{ title }}
    <svg [attr.viewBox]="viewBox()">
      <g *ngFor="let chunk of chunks()">
        <title>{{ jsonDump(chunk) }}</title>
        <rect [attr.x]="chunk.cx" [attr.y]="chunk.cz" width="0.8" height="0.8"></rect>
      </g>
    </svg>
    <router-outlet></router-outlet>
  `,
  styles: []
})
export class AppComponent implements OnInit {
  title = 'bedrock-ui';
  chunks = () => this.gameDataService.overworldChunks().map(e => this.gameDataService.chunkKeyInfo(`${e}32`));
  jsonDump = data => JSON.stringify(data, null, 2);
  viewBox = () => {
    let minX = Math.min(...this.chunks().map(e => e['cx']));
    let minY = Math.min(...this.chunks().map(e => e['cz']));
    let width = Math.max(...this.chunks().map(e => e['cx'])) - minX;
    let height = Math.max(...this.chunks().map(e => e['cz'])) - minY;
    let out = `${minX} ${minY} ${width} ${height}`;
    // return some default if any values are Infinity or NaN
    return /^[\d -]+$/.test(out) ? out : "-16 -16 32 32";
  }
  constructor(private gameDataService: GameDataService) { }
  ngOnInit() {
    this.gameDataService.refreshKeys();
    let myfoo = '0100000001000000';
    console.log(this.gameDataService.chunkKeyInfo(`${myfoo}32`));
  }

}
