import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AccountComponent } from './account.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgZorroAntdModule, NZ_I18N, zh_CN } from 'ng-zorro-antd';
import { LoggerService } from '../commons/services/logger.service';
import { AccountgroupComponent } from './group/account-group.component';
import { PolicyComponent } from './policy/policy.component';
import { TenantComponent } from './tenant/tenant.component';
import { SelectuserModule } from './../commons/selectuser/selectuser.module';

@NgModule({
  imports: [
    CommonModule,
    NgZorroAntdModule,
    FormsModule,
    ReactiveFormsModule,
    SelectuserModule
  ],
  providers: [{ provide: NZ_I18N, useValue: zh_CN }, LoggerService],
  declarations: [AccountComponent, AccountgroupComponent, PolicyComponent, TenantComponent]
})
export class AccountModule { }
