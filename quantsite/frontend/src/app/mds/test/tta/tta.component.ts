import { Component, OnInit } from '@angular/core';
import { BreadcrumbService } from '../breadcrumb/breadcrumb.service';

@Component({
  selector: 'app-tta',
  templateUrl: './tta.component.html',
  styleUrls: ['./tta.component.scss']
})
export class TtaComponent implements OnInit {

  constructor(
    private breadcrumbService: BreadcrumbService
  ) { }

  ngOnInit() {
    this.breadcrumbService.freshBreadcrumbs();
  }

}
