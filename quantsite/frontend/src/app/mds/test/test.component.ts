import { Component, OnInit } from '@angular/core';
import { BreadcrumbService } from './breadcrumb/breadcrumb.service';

@Component({
  selector: 'app-test',
  templateUrl: './test.component.html',
  styleUrls: ['./test.component.scss']
})
export class TestComponent implements OnInit {

  constructor(
    private breadcrumbService: BreadcrumbService
  ) {
  }

  ngOnInit() {
    this.breadcrumbService.freshBreadcrumbs();
  }

}
