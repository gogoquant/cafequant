import { Component, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  FormControl,
  Validators
} from '@angular/forms';
import { TenantService } from './tenant.service';
import { LoggerService } from 'src/app/commons/services/logger.service';
import { NzMessageService, NzModalRef, NzModalService } from 'ng-zorro-antd';

@Component({
  selector: 'app-tenant',
  templateUrl: './tenant.component.html',
  styleUrls: ['./tenant.component.scss']
})
export class TenantComponent implements OnInit {

  delConfirmModal: NzModalRef;

  /** 是否显示租户列表页面 */
  showTenant = true;
  pageNum = 1;
  pageSize = 10;
  total = 1;
  loading = false;
  keyword = '';
  /** 查询列表结果 */
  groupData = [];

  /** 是否显示新增租户页面 */
  showAddTenant = false;
  validateForm: FormGroup;
  /** 租户管理员 */
  userName = '';

  constructor(
    private tenantService: TenantService,
    private logger: LoggerService,
    private fb: FormBuilder,
    private nzMessageService: NzMessageService,
    private modal: NzModalService
  ) { }

  ngOnInit() {
    /** 进入页面查询默认所有分页租户 */
    this.searchData();
    /** 注册校验表单 */
    this.validateForm = this.fb.group({
      tenant_id: [null, [Validators.required, Validators.maxLength, this.tenantValidator]],
      name: [null, [Validators.required, Validators.maxLength]],
      manager: [null, [Validators.required]],
      comment: [null, [Validators.maxLength]]
    });
  }
  /** 调用选择用户公共组件赋值给RTC账户名，userData 是个数组 */
  handleUserOk(userData) {
    this.validateForm.get('manager').setValue(userData[0].username);
  }

  /**根据分页查询租户列表 */
  queryTenants(): void {
    this.tenantService.queryTenants(this.pageNum, this.pageSize, this.keyword)
      .subscribe((data: any) => {
        this.groupData = data.data;
        this.pageNum = data.page_index;
        this.pageSize = data.page_size;
        this.total = data.total;
        this.loading = false;
      });
  }

  /** 分页查询用，根据reset是否选择了pageNum数做条件 */
  searchData(reset: boolean = false): void {
    if (reset) {
      this.pageNum = 1;
    }
    this.loading = true;
    this.queryTenants();
  }

  /**根据输入的内容查询租户列表 */
  queryTenantByName(event: any): void {
    const tenantName = event.target.value;
    this.keyword = tenantName;
    this.logger.info(this.keyword);
    this.queryTenants();
  }

  /** 显示新增租户页面 */
  showAddTenantView(): void {
    this.showAddTenant = true;
    this.showTenant = false;
  }

  /** 根据正则表达式校验规则 */
  tenantValidator = (control: FormControl): { [s: string]: boolean } => {
    // const reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/;
    const reg = /^[a-z]([-a-z0-9]{0,13}[a-z0-9])?$/;
    if (!control.value) {
      return { required: true };
    } else if (!reg.test(control.value)) {
      return { tenant_id: true, error: true };
    }
  }

  submitForm = ($event, value) => {
    $event.preventDefault();
    for (const key of Object.keys(this.validateForm.controls)) {
      this.validateForm.controls[key].markAsDirty();
      this.validateForm.controls[key].updateValueAndValidity();
    }
  }
  resetForm(e: MouseEvent): void {
    e.preventDefault();
    this.validateForm.reset();
    for (const key of Object.keys(this.validateForm.controls)) {
      this.validateForm.controls[key].markAsDirty();
      this.validateForm.controls[key].updateValueAndValidity();
    }
  }

  /** 新增租户提交数据到接口 */
  addTenant($event): void {
    if (this.validateForm.valid) {
      const params = {
        tenant_id: this.validateForm.value.tenant_id,
        name: this.validateForm.value.name,
        manager: this.validateForm.value.manager,
        comment: this.validateForm.value.comment,
      };
      this.tenantService.AddTenant(params).subscribe((data: any) => {
        this.nzMessageService.success('添加租户成功！');
        this.searchData();
        this.clickList();
      });
    }
  }

  /** 返回列表事件 */
  clickList(): void {
    this.showAddTenant = false;
    this.showTenant = true;
  }

  /** 根据租户UUID删除租户 */
  delTenantConfirm(name: string, uuid: string): void {
    this.logger.info(uuid);
    this.delConfirmModal = this.modal.confirm({
      nzTitle: '您确定要删除' + name,
      nzContent: '点确定删除该租户记录，点取消则不删除！',
      nzOnOk: () => this.tenantService.deleteTenant(uuid).subscribe((data: any) => {
        this.searchData();
      }, (error) => {
        this.nzMessageService.error(error.error);
      })
    });
  }

}
