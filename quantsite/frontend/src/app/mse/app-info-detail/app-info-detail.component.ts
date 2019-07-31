import { Component, OnInit } from '@angular/core';

import { AppInfo } from '../beans';
import { MseService } from '../mse.service';
import { ActivatedRoute, Router, ParamMap } from '@angular/router';
import { switchMap } from 'rxjs/operators';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-app-info-detail',
  templateUrl: './app-info-detail.component.html',
  styleUrls: ['./app-info-detail.component.scss']
})
export class AppInfoDetailComponent implements OnInit {
  private currentAppInfo$: Observable<AppInfo>;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private mseService: MseService
  ) { }

  ngOnInit() {
    this.currentAppInfo$ = this.route.paramMap.pipe(
      switchMap((params: ParamMap) =>
        this.mseService.queryAppInfo(Number.parseInt(params.get('id')))));
  }
}
