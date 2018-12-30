import { Component, OnInit } from '@angular/core';

import { GameDataService } from './game-data.service';

@Component({
  selector: 'app-root',
  template: `
    Under Construction
    <li *ngFor="let foo of gameDataService.keyList">{{ foo }}</li>
    <router-outlet></router-outlet>
  `,
  styles: []
})
export class AppComponent implements OnInit {
  title = 'bedrock-ui';
  constructor(private gameDataService: GameDataService) { }
  ngOnInit() {
    this.gameDataService.refreshKeys();
  }

}
