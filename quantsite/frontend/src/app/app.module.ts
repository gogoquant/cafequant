import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { NgZorroAntdModule, NZ_I18N, zh_CN } from 'ng-zorro-antd';
import { registerLocaleData } from '@angular/common';
import zh from '@angular/common/locales/zh';
import { AppRoutingModule } from './app-routing.module';
import { NgxEchartsModule } from 'ngx-echarts';

import { AccountModule } from './account/account.module';
import { LoggingModule } from './logging/logging.module';
import { MceModule } from './mce/mce.module';
import { SystemModule } from './system/system.module';
import { MdsModule } from './mds/mds.module';
// import { MseModule } from './mse/mse.module';

// Ngx
import { NgxPermissionsModule } from 'ngx-permissions';


// service
import { LoggerService } from './commons/services/logger.service';


registerLocaleData(zh);

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    FormsModule,
    HttpClientModule,
    NgxEchartsModule,
    NgZorroAntdModule,
    MceModule,
    MdsModule,
    NgxPermissionsModule.forRoot(),
    AccountModule,
    LoggingModule,
    SystemModule,
    AppRoutingModule,
  ],
  providers: [{ provide: NZ_I18N, useValue: zh_CN }, LoggerService],
  bootstrap: [AppComponent]
})
export class AppModule { }
