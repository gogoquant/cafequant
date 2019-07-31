import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LoggingComponent } from './logging.component';
// import { NgZorroAntdModule } from '../../../node_modules/ng-zorro-antd';
import { NgZorroAntdModule, NZ_I18N, zh_CN } from 'ng-zorro-antd';
import { LoggerService } from '../commons/services/logger.service';
import { AppRoutingModule } from '../app-routing.module';
@NgModule({
  imports: [
    CommonModule,
    NgZorroAntdModule
  ],
  providers: [{ provide: NZ_I18N, useValue: zh_CN }, LoggerService],
  declarations: [LoggingComponent]
})
export class LoggingModule { }
