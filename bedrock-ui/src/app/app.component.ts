import { Component, OnInit } from '@angular/core';
import { ChunkListService } from './chunk-list.service';

@Component({
  selector: 'app-root',
  template: `
    Under Construction
    <router-outlet></router-outlet>
  `,
  styles: []
})
export class AppComponent implements OnInit {
  title = 'bedrock-ui';
  constructor(private chunkListService: ChunkListService) { }
  ngOnInit() {
    this.chunkListService.RefreshKeys();
  }

}
