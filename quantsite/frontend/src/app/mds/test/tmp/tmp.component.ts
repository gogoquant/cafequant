import { Component, OnInit } from '@angular/core';
import { BreadcrumbService } from '../breadcrumb/breadcrumb.service';

@Component({
  selector: 'app-tmp',
  templateUrl: './tmp.component.html',
  styleUrls: ['./tmp.component.scss']
})
export class TmpComponent implements OnInit {

  constructor(
    private breadcrumbService: BreadcrumbService
  ) { }

  ngOnInit() {
    this.breadcrumbService.freshBreadcrumbs();
  }

}
