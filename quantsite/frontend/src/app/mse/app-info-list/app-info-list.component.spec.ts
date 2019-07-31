import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AppInfoListComponent } from './app-info-list.component';

describe('AppInfoListComponent', () => {
  let component: AppInfoListComponent;
  let fixture: ComponentFixture<AppInfoListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AppInfoListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AppInfoListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
