import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AgileProjectComponent } from './agile-project.component';

describe('AgileProjectComponent', () => {
  let component: AgileProjectComponent;
  let fixture: ComponentFixture<AgileProjectComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AgileProjectComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AgileProjectComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
