import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { NgZorroAntdModule, NZ_I18N, zh_CN } from 'ng-zorro-antd';
import { LoggerService } from '../commons/services/logger.service';
import { MceComponent } from './mce.component';
import { ClusterComponent } from './cluster/cluster.component';

@NgModule({
    imports: [
        CommonModule,
        NgZorroAntdModule,
        FormsModule
    ],
    providers: [{ provide: NZ_I18N, useValue: zh_CN }, LoggerService],
    declarations: [MceComponent, ClusterComponent]
})
export class MceModule { }
