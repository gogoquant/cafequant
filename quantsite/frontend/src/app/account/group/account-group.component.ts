import { Component, OnInit } from '@angular/core';
import { AccountGroupService } from './account-group.service';
import { LoggerService } from './../../commons/services/logger.service';
import { NzMessageService } from 'ng-zorro-antd';
import {
  FormBuilder,
  FormControl,
  FormGroup,
  Validators
} from '@angular/forms';

@Component({
  selector: 'app-account-group',
  templateUrl: './account-group.component.html',
  styleUrls: ['./account-group.component.scss']
})
export class AccountgroupComponent implements OnInit {
  modulesData = [];
  prSelectedValue = '';
  grSelectedValue = 'DEV';
  urSelectedValue = '管理员';
  groupName = '';
  showGroup = true;
  isVisible = false;
  isEditVisible = false;
  isUserVisible = false;
  groupData = [];
  userData = [];
  validateForm: FormGroup;
  validateEditForm: FormGroup;
  validateUserForm: FormGroup;
  editRowId: 0;
  list = [];
  rightData = [];
  leftData = [];
  checkedGroupId: number;

  constructor(
    private accountGroupService: AccountGroupService,
    private logger: LoggerService,
    private fb: FormBuilder,
    private nzMessageService: NzMessageService
  ) { }

  ngOnInit() {
    this.getModules();
    this.validateForm = this.fb.group({
      group_name       : [ null, [ Validators.required ] ],
      group_id         : [ null, [ Validators.required ] ],
      project_id       : [ null, [ Validators.required ] ],
      description      : [ null ]
    });
    this.validateEditForm = this.fb.group({
      group_name       : [ null, [ Validators.required ] ],
      group_id         : [ null, [ Validators.required ] ],
      project_id       : [ null, [ Validators.required ] ],
      description      : [ null ]
    });
    this.validateUserForm = this.fb.group({
      user_id       : [ null ]
    });
  }
  getModules(): void {
    this.accountGroupService.getModules().subscribe((data: any) => {
        this.modulesData = data;
        this.prSelectedValue = data[0].ID;
        this.getGroupsById(this.prSelectedValue);
        console.log('modulesData', this.modulesData);
      });
  }
  getGroupsById = (id) => {
    this.accountGroupService.getGroupsById(id).subscribe((data: any) => {
      this.groupData = data;
    });
  }
  projectChange = (value) => {
    this.getGroupsById(value);
  }
  groupClick = (name, id) => {
    this.groupName = name;
    this.showGroup = false;
    this.getUserByGroup(id);
    this.checkedGroupId = id;
  }
  userClick = () => {
    this.showGroup = true;
  }
  showAddModal = () => {
    this.isVisible = true;
  }
  handleOk = () => {
    this.submitForm();
  }
  handleCancel = () => {
    this.isVisible = false;
  }
  submitForm = () => {
    for (const i of Object.keys(this.validateForm.controls)) {
      this.validateForm.controls[ i ].markAsDirty();
      this.validateForm.controls[ i ].updateValueAndValidity();
    }
    if (this.validateForm.valid) {
      const params = {
        group_name: this.validateForm.value.group_name,
        group_id: this.validateForm.value.group_id,
        project_id: this.validateForm.value.project_id,
        description: this.validateForm.value.description,
      };
      this.accountGroupService.AddGroup(params).subscribe((data: any) => {
        this.nzMessageService.success('添加用户组成功');
        this.getGroupsById(this.prSelectedValue);
        this.isVisible = false;
      });
    }
  }
  delGroupCancel = () => {
  }
  delGroupConfirm = (id) => {
    this.accountGroupService.DelGroup(id).subscribe((data: any) => {
      this.getGroupsById(this.prSelectedValue);
      this.nzMessageService.success('删除用户组成功');
    });
  }
  showEditModal = (data) => {
    this.validateEditForm.get('group_name').setValue(data.group_name);
    this.validateEditForm.get('group_id').setValue(data.group_id);
    this.validateEditForm.get('project_id').setValue(data.project_id);
    this.validateEditForm.get('description').setValue(data.description);
    this.editRowId = data.id;
    this.isEditVisible = true;
  }
  handleEditOk = () => {
    for (const i of Object.keys(this.validateEditForm.controls)) {
      this.validateEditForm.controls[ i ].markAsDirty();
      this.validateEditForm.controls[ i ].updateValueAndValidity();
    }
    if (this.validateEditForm.valid) {
      const params = {
        id: this.editRowId,
        group_name: this.validateEditForm.value.group_name,
        group_id: this.validateEditForm.value.group_id,
        project_id: this.validateEditForm.value.project_id,
        description: this.validateEditForm.value.description,
      };
      this.accountGroupService.EditGroup(params).subscribe((data: any) => {
        this.nzMessageService.success('编辑用户组成功');
        this.getGroupsById(this.prSelectedValue);
        this.isEditVisible = false;
      });
    }
  }
  handleEditCancel = () => {
    this.isEditVisible = false;
  }
  getUserByGroup = (id) => {
    this.accountGroupService.getUserByGroup(id).subscribe((data: any) => {
      this.userData = data;
    });
  }
  delUserCancel = () => {
  }
  delUserConfirm = (id) => {
    this.accountGroupService.DelUser(id).subscribe((data: any) => {
      this.getUserByGroup(this.checkedGroupId);
      this.nzMessageService.success('删除用户成功');
    });
  }
  showUserModal = () => {
    this.isUserVisible = true;
  }
  handleUserOk = () => {
    const params = [];
    this.rightData.forEach(item => {
      const listitem = {
        group_id: this.checkedGroupId,
        user_id: item.user_id
      };
      params.push(listitem);
    });
    this.accountGroupService.AddUsers(params).subscribe((data: any) => {
      this.nzMessageService.success('添加用户成功');
      this.getUserByGroup(this.checkedGroupId);
      this.isUserVisible = false;
    });
  }
  handleUserCancel = () => {
    this.isUserVisible = false;
  }
  getData = (value) => {
    this.accountGroupService.getAccounts(value).subscribe((data: any) => {
      if (data.data.length > 0) {
        data.data.forEach((element, index) => {
          element.key = element.user_id;
          element.title = element.name;
          element.direction = '';
        });
      }
      this.leftData = data.data;
      this.list = this.leftData.concat(this.rightData);
    });
  }
  reload(direction: string): void {
    this.getData('');
    this.nzMessageService.success(`your clicked ${direction}!`);
  }
  search = (ret) => {
    if (ret.direction === 'left') {
      if (ret.value === '') {
        this.leftData = [];
        this.list = this.leftData.concat(this.rightData);
      } else {
        this.getData(ret.value);
      }
    }
  }
  select(ret): void {
    // console.log('nzSelectChange', ret);
  }
  change = (ret) => {
    if (ret.from === 'left') {
      this.rightData = this.rightData.concat(ret.list);
      this.rightData.forEach(item => {
        item.direction = 'right';
      });
    }
    if (ret.from === 'right') {
      const actionData = this.rightData;
      ret.list.forEach(item => {
        this.rightData.forEach((element, index) => {
          if (element.user_id === item.user_id) {
            actionData.splice(index, 1);
          }
        });
      });
      this.rightData = actionData;
    }
  }
}
