import { Component, OnInit, OnChanges, Input } from '@angular/core';
import { IBreadcrumb } from './breadcrumb.model';
import { BreadcrumbService } from './breadcrumb.service';
import { Subject, Observable } from 'rxjs';

@Component({
  selector: 'app-breadcrumb',
  templateUrl: './breadcrumb.component.html',
  styleUrls: ['./breadcrumb.component.scss']
})
export class BreadcrumbComponent implements OnInit, OnChanges {

  breadcrumbObservable: Observable<IBreadcrumb[]>;
  @Input() prefixUrl: string;

  constructor(
    private breadCrumbeService: BreadcrumbService
  ) {}

  ngOnInit() {
    setTimeout(() => {
      this.breadcrumbObservable = this.breadCrumbeService.breadcrumbSubject.asObservable();
      this.breadCrumbeService.setPrefixUrl(this.prefixUrl);
    }, 0);
    setTimeout(() => {
      this.breadCrumbeService.freshBreadcrumbs();
    });
  }

  ngOnChanges() {
    this.breadCrumbeService.setPrefixUrl(this.prefixUrl);
  }

}
