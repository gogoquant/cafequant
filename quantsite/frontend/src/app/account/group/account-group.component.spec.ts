import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AccountgroupComponent } from './account-group.component';

describe('AccountgroupComponent', () => {
  let component: AccountgroupComponent;
  let fixture: ComponentFixture<AccountgroupComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AccountgroupComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AccountgroupComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
