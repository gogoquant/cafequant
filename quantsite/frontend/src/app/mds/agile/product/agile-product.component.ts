import { Component, OnInit } from '@angular/core';
import { AccountService } from '../../../account/account.service';
import { LoggerService } from '../../../commons/services/logger.service';
import { AgileService } from '../agile.service';

@Component({
  selector: 'app-agile-product',
  templateUrl: './agile-product.component.html',
  styleUrls: ['./agile-product.component.scss']
})
export class AgileProductComponent implements OnInit {

  /** 调用的账户获取信息 */
  name = '';
  user_id = '';
  /**当前登陆用户名 */
  currUserName = 'wei9.li';
  productName = ``;
  pageNum = 1;
  pageSize = 10;
  total = 1;
  loading = false;
  /** 租户id */
  tenantId = '';

  showProduct = true;
  /**查询的列表结果 */
  groupData = [];

  constructor(private accountService: AccountService,
    private logger: LoggerService,
    private agileService: AgileService) { }


  ngOnInit() {
    this.searchData();
  }

  /** 调用选择用户公共组件赋值给RTC账户名，userData 是个数组 */
  handleUserOk(userData) {
    this.name = userData[0].username;
    this.accountService.getUserIdByUserName(this.name).subscribe((data: any) => {
      this.name = data.name;
      this.user_id = data.user_id;
      this.queryProducts(this.user_id, this.productName, this.tenantId);
    });
  }

  /**根据输入的RTC账户名获取user对象信息并赋值userId用于查询创建者的产品列表 */
  onblur(event: any): void {
    const userName = event.target.value;
    if (userName !== '' && userName !== undefined) {
      this.accountService.getUserIdByUserName(userName).subscribe((data: any) => {
        this.name = data.name;
        this.user_id = data.user_id;
        this.queryProducts(this.user_id, this.productName, this.tenantId);
      });
    } else {
      this.name = '';
      this.user_id = '';
    }
  }

  /**根据输入的RTC账户名获取userId或输入的产品名称查询产品列表 */
  queryProducts(creator: string, productName: string, tenantId: string): void {
    if (this.currUserName !== '' && this.currUserName !== undefined) {
      //      this.logger.info(this.productName);
      this.agileService.getProducts(this.pageNum, this.pageSize, this.currUserName, creator, productName, tenantId)
        .subscribe((data: any) => {
          this.groupData = data.list;
          this.pageNum = data.pageNum;
          this.pageSize = data.pageSize;
          this.total = data.total;
          this.loading = false;
          this.logger.info(this.total);

        });
    }
  }

  /** 产品名称输入后的查询 */
  productNameBlur(event: any): void {
    const value = event.target.value;
    if (value !== '' && value !== undefined) {
      this.productName = value;
      this.queryProducts(this.user_id, this.productName, this.tenantId);
    } else {
      this.productName = '';
    }
  }

  /** 分页查询用 */
  searchData(reset: boolean = false): void {
    if (reset) {
      this.pageNum = 1;
    }
    this.loading = true;
    this.queryProducts(this.user_id, this.productName, this.tenantId);
  }

}
