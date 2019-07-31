import { CommonModule } from '@angular/common';
import { DashboardComponent } from './dashboard/dashboard.component';
import { NgModule } from '@angular/core';
import { NgxEchartsModule } from 'ngx-echarts';
@NgModule({
  imports: [
    CommonModule, NgxEchartsModule
  ],
  declarations: [DashboardComponent],
  schemas: [
    // CUSTOM_ELEMENTS_SCHEMA,
  ]
})
export class SystemModule { }
