import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MseComponent } from './mse.component';
import { FormsModule } from '@angular/forms';

import { NgZorroAntdModule, NZ_I18N, zh_CN } from 'ng-zorro-antd';
import { LoggerService } from '../commons/services/logger.service';
import { AppInfoDetailComponent } from './app-info-detail/app-info-detail.component';
import { AppInfoListComponent } from './app-info-list/app-info-list.component';
import { BreadcrumbComponent } from './breadcrumb/breadcrumb.component';
import { MseRoutingModule } from './mse-routing.module';
import { MseService } from './mse.service';


@NgModule({
  imports: [
    CommonModule, NgZorroAntdModule, FormsModule, MseRoutingModule
  ],
  providers: [{ provide: NZ_I18N, useValue: zh_CN }, LoggerService, MseService],
  declarations: [MseComponent, AppInfoDetailComponent, AppInfoListComponent, BreadcrumbComponent]
})
export class MseModule { }
