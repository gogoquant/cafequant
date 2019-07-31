import { Component, OnInit, TemplateRef } from '@angular/core';
import { AccountService } from '../../../account/account.service';
import { LoggerService } from '../../../commons/services/logger.service';
import { AgileService } from '../agile.service';
import {
  FormBuilder,
  FormGroup,
  Validators
} from '@angular/forms';
import { NzModalService, NzModalRef, NzMessageService } from 'ng-zorro-antd';
import { map } from 'rxjs/operators';

@Component({
  selector: 'app-agile-project',
  templateUrl: './agile-project.component.html',
  styleUrls: ['./agile-project.component.scss']
})
export class AgileProjectComponent implements OnInit {

  /** 调用的账户获取信息 */
  name = '';
  user_id = '';
  /**当前登陆用户名 */
  currUserName = 'wei9.li';
  productName = ``;
  projectName = ``;
  pageNum = 1;
  pageSize = 10;
  total = 1;
  loading = false;
  /** 租户id */
  tenantId = '';
  /** 是否显示项目列表 */
  showProject = true;
  /**查询的列表结果 */
  groupData = [];

  /** 是否显示项目详情 */
  showDetail = false;
  projectDetail = {};
  validateForm: FormGroup;

  /** 项目配置详情 */
  projectConfig = {};
  /** 项目配置modal */
  configModal: NzModalRef;
  tplModalButtonLoading = false;
  htmlModalVisible = false;
  projectConfigForm: FormGroup;

  /** 项目成员数据列表 */
  projectMembers = {};
  /** 角色map */
  map: { [key: string]: string } = {};

  constructor(private accountService: AccountService,
    private logger: LoggerService,
    private agileService: AgileService,
    private fb: FormBuilder,
    private modalService: NzModalService,
    private message: NzMessageService
  ) { }

  ngOnInit() {
    this.searchData();
    this.validateForm = this.fb.group({
      startTime: [null],
      dateRange: [null]
    });
    this.projectConfigForm = this.fb.group({
      updateTime: [null]
    });
  }

  /** 调用选择用户公共组件赋值给RTC账户名 */
  handleUserOk(userData) {
    this.name = userData[0].username;
    this.accountService.getUserIdByUserName(this.name).subscribe((data: any) => {
      this.name = data.name;
      this.user_id = data.user_id;
      this.queryProjects(this.user_id, this.productName, this.projectName, this.tenantId);
    });
  }


  /**根据输入的RTC账户名获取user对象信息并赋值userId用于查询创建者的项目列表 */
  onblur(event: any): void {
    const userName = event.target.value;
    if (userName !== '' && userName !== undefined) {
      this.accountService.getUserIdByUserName(userName).subscribe((data: any) => {
        this.name = data.name;
        this.user_id = data.user_id;
        this.queryProjects(this.user_id, this.productName, this.projectName, this.tenantId);
      });
    } else {
      this.name = '';
      this.user_id = '';
    }
  }

  /**根据输入的RTC账户名获取userId或输入的产品名称、项目名称查询项目列表 */
  queryProjects(creator: string, productName: string, projectName: string, tenantId: string): void {
    if (this.currUserName !== '' && this.currUserName !== undefined) {
      //      this.logger.info(this.productName);
      this.agileService.getProjects(this.pageNum, this.pageSize, this.currUserName, creator, productName, projectName, tenantId)
        .subscribe((data: any) => {
          this.groupData = data.list;
          this.pageNum = data.pageNum;
          this.pageSize = data.pageSize;
          this.total = data.total;
          this.loading = false;
          // this.logger.info(this.total);
        });
    }
  }

  /** 产品名称输入后的查询 */
  productNameBlur(event: any): void {
    const value = event.target.value;
    if (value !== '' && value !== undefined) {
      this.productName = value;
      this.queryProjects(this.user_id, this.productName, this.projectName, this.tenantId);
    } else {
      this.productName = '';
    }
  }

  /** 项目名称输入后的查询 */
  projectNameBlur(event: any): void {
    const value = event.target.value;
    if (value !== '' && value !== undefined) {
      this.projectName = value;
      this.queryProjects(this.user_id, this.productName, this.projectName, this.tenantId);
    } else {
      this.projectName = '';
    }
  }

  /** 分页查询用 */
  searchData(reset: boolean = false): void {
    if (reset) {
      this.pageNum = 1;
    }
    this.loading = true;
    this.queryProjects(this.user_id, this.productName, this.projectName, this.tenantId);
  }


  /**根据projectId查看项目详情 */
  getProjectDetail(projectId: string): void {
    // this.logger.info(projectId);
    this.showDetail = true;
    this.showProject = false;
    this.agileService.getProjectById(projectId, this.currUserName).subscribe((data: any) => {
      this.projectDetail = data;
      // this.validateForm.get('startTime').setValue(data.startTime);
      this.validateForm.get('dateRange').setValue([data.startTime, data.endTime]);
      // this.logger.info(this.validateForm.get('dateRange').value);
    });
  }

  /** 返回列表事件 */
  clickList(): void {
    this.showProject = true;
    this.showDetail = false;
  }

  /** 根据项目编号id查询项目配置信息 */
  getProjectconfig(projectId: string, callback): boolean {
    let flag = false;
    if (projectId !== '' && projectId !== undefined) {
      this.agileService.getProjectConfig(projectId, this.currUserName).subscribe((data: any) => {
        this.projectConfig = data;
        this.projectConfigForm.get('updateTime').setValue(data.updateTime);
        // projectId存在，则配置信息存在且正确返回OK
        if (data.projectId) {
          flag = true;
        }
        this.logger.info(this.projectConfig);
        callback(flag);
      });
    }
    return flag;
  }

  /** 项目配置链接生成modal框并展示信息 */
  projectConfigModal(projectName: string, projectId: string, tplContent: TemplateRef<{}>): void {
    this.getProjectconfig(projectId, (flag: boolean) => {
      if (flag) {
        const modal = this.modalService.create({
          nzTitle: `${projectName}项目配置信息`,
          nzContent: tplContent,
          // 'pass array of button config to nzFooter to create multiple buttons',
          nzFooter: [
            {
              label: '关闭',
              shape: 'default',
              onClick: () => {
                modal.destroy();
              }
            }
          ],
          nzComponentParams: this.projectConfig,
          nzWidth: 780,
          nzZIndex: 1001
        });
      } else {
        // 类似alert的功能
        this.message.info('暂无项目配置信息！');
      }
    });
  }

  /** 获取项目成员列表信息 */
  getProjectMembers(projectId: string, callback): boolean {
    let flag = false;
    if (projectId !== '' && projectId !== undefined) {
      this.agileService.getProjectUsers(projectId, this.currUserName).subscribe((data: any) => {
        this.projectMembers = data;
        if (data.roleDTOList) {
          flag = true;
          const roles: any = data.roleDTOList;
          // 遍历角色记录列表并保存到map对象中，与下边map取值成对使用；
          // 相对于从roleList 和 roleUserRelationDTOList 嵌套循环比较赋值map方式相对于从roleList只循环一次
          // for (const role of roles) {
          //   this.map[role.roleId] = role.roleName;
          // }
        }
        callback(flag);
        // 也可以直接遍历每个成员记录，将map中的角色赋值，在列表中遍历的时候直接取值即可
        // data.roleUserRelationDTOList.forEach(element => {
        //   element.roleName = this.map[element.roleId];
        // });
      });
    }
    return flag;
  }

  /** 项目成员生成modal框并展示信息 */
  projectMemberModal(projectName: string, projectId: string, tplContent: TemplateRef<{}>): void {
    this.getProjectMembers(projectId, (flag: boolean) => {
      if (this.projectMembers !== {} && this.projectMembers !== undefined) {
        const modal = this.modalService.create({
          nzTitle: `${projectName}项目成员信息`,
          nzContent: tplContent,
          // 'pass array of button config to nzFooter to create multiple buttons',
          nzFooter: [
            {
              label: '关闭',
              shape: 'default',
              onClick: () => modal.destroy()
            }
          ],
          nzComponentParams: this.projectConfig,
          nzWidth: 660,
          nzZIndex: 1001
        });
      } else {
        this.message.info('暂无项目成员信息！');
      }
    });

  }



}
