import { Component, OnInit, EventEmitter, Output } from '@angular/core';
import { SelectuserService } from './selectuser.service';
import {
  FormBuilder,
  FormControl,
  FormGroup,
  Validators
} from '@angular/forms';

@Component({
  selector: 'app-selectuser',
  templateUrl: './selectuser.component.html',
  styleUrls: ['./selectuser.component.scss']
})
export class SelectuserComponent implements OnInit {
  isUserVisible = false;
  list = [];
  rightData = [];
  leftData = [];
  validateUserForm: FormGroup;
  @Output() handleOk = new EventEmitter<any>();
  constructor(
    private selectuserService: SelectuserService,
    private fb: FormBuilder,
  ) { }

  ngOnInit() {
    this.validateUserForm = this.fb.group({
      user_id       : [ null ]
    });
  }

  showUserModal = () => {
    this.isUserVisible = true;
  }
  handleUserCancel = () => {
    this.isUserVisible = false;
  }
  handleUserOk = () => {
    const userData = [];
    this.rightData.forEach(item => {
      const listitem = {
        username: item.username,
        user_id: item.user_id,
        name: item.name,
        email: item.email,
      };
      userData.push(listitem);
    });
    this.handleOk.emit(userData);
    this.isUserVisible = false;
  }
  getData = (value) => {
    this.selectuserService.getAccounts(value).subscribe((data: any) => {
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
