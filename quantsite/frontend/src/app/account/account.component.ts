import { Component, OnInit } from '@angular/core';
import { AccountService } from './account.service';
import { NzMessageService } from 'ng-zorro-antd';
import { LoggerService } from './../commons/services/logger.service';

@Component({
  selector: 'app-account',
  templateUrl: './account.component.html',
  styleUrls: ['./account.component.scss']
})
export class AccountComponent implements OnInit {
  // sync user loading
  syncUserLoading = false;

  accounts = [];
  pageIndex = 1;
  pageSize = 10;
  total = 1;
  loading = true;
  limit = '';
  keywords = '';

  constructor(private accountService: AccountService,
    private logger: LoggerService,
    private message: NzMessageService) { }

  ngOnInit() {
    this.searchData();
  }

  searchData(reset: boolean = false): void {
    if (reset) {
      this.pageIndex = 1;
    }
    this.loading = true;
    this.accountService.getAccounts(this.keywords, this.limit,
      this.pageIndex, this.pageSize).subscribe((data: any) => {
        this.loading = false;
        this.total = data.total;
        this.accounts = data.data;
      });
  }
  handleUserOk(userData) {
    console.log('userData', userData);
  }
  // sync users from ldap
  syncUsers(): void {
    this.syncUserLoading = true;
    this.accountService.syncUsers().subscribe((data: any) => {
      this.logger.info(data);
      this.syncUserLoading = false;
      if (data != null) {
        this.message.create('success', `同步用户数${data.insert},总用户数:${data.total}`);
      }
    });
  }
}
