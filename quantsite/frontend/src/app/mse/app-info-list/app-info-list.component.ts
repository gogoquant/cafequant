import { Component, OnInit } from '@angular/core';
import { Subject, Observable, of } from 'rxjs';
import { debounceTime, distinctUntilChanged, switchMap } from 'rxjs/operators';

import { MseService } from '../mse.service';
import { AppInfo } from '../beans';

class SearchCondition {
  constructor(public orgName: string, public owner: string) {
  }
}

@Component({
  selector: 'app-app-info-list',
  templateUrl: './app-info-list.component.html',
  styleUrls: ['./app-info-list.component.scss']
})
export class AppInfoListComponent implements OnInit {
  private currentSearchCondition: SearchCondition;
  private appInfoList: AppInfo[];
  private currentPageData: AppInfo[];
  private orgNames: string[];
  private selectedDepartment: string;

  // table varables.
  private isLoading: boolean;
  private total: number;
  private pageSize: number;
  private pageIndex: number;

  private isSearchConditionValueChanged: boolean;

  private keywordText$ = new Subject<string>();
  private keywordObserableHolder$: Observable<{}>;

  constructor(
    public mseService: MseService
  ) { }

  ngOnInit() {
    this.isSearchConditionValueChanged = true;
    this.orgNames = [];
    this.pageIndex = 1;
    this.pageSize = 10;
    this.currentSearchCondition = new SearchCondition(null, null);
    this.getOrgNames();
    this.initObserable();

    console.log('organizations.length=');
    setTimeout(() => this.refreshTableData(), 500);
  }

  private initObserable(): void {
    this.keywordObserableHolder$ = this.keywordText$
    .pipe(
      debounceTime(500),
      distinctUntilChanged(),
      switchMap((keyword) => {
        console.log('keyword search tigered.');
        this.isSearchConditionValueChanged = true;
        this.currentSearchCondition.owner = keyword;
        setTimeout(() => this.refreshTableData(), 0);
        return of(); })
    );
  }

  private onSearchBtnClicked(): void {
    this.isSearchConditionValueChanged = true;
    setTimeout(() => this.refreshTableData(), 0);
  }

  private getOrgNames(): void {
    this.mseService.queryOrganizations()
      .subscribe(organizations => {
        console.log('organizations.length=', organizations.length);
        organizations.forEach(org => this.orgNames.push(org.orgName));
      });
  }

  private refreshTableData(): void {
    if (true === this.isSearchConditionValueChanged) {
      this.isLoading = true;
      this.mseService.queryApps(this.currentSearchCondition.orgName, this.currentSearchCondition.owner, 0, 1000)
          .subscribe(_appInfoList => {
            this.total = _appInfoList.length;
            this.appInfoList = _appInfoList;
            this.currentPageData = this.appInfoList.slice((this.pageIndex - 1) * this.pageSize, (this.pageIndex) * this.pageSize);
            this.isSearchConditionValueChanged = false;
            this.isLoading = false;
          });
    } else {
      this.currentPageData = this.appInfoList.slice((this.pageIndex - 1) * this.pageSize, (this.pageIndex) * this.pageSize);
    }
  }

  private pageIndexchanged(): void {
    this.refreshTableData();
  }

  private pageSizeChange(): void {
    this.refreshTableData();
  }

  private onOwnerChanged(keyword: string) {
    this.keywordText$.next(keyword);
  }

}
