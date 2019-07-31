import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AgileProductComponent } from './agile-product.component';

describe('AgileProductComponent', () => {
  let component: AgileProductComponent;
  let fixture: ComponentFixture<AgileProductComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AgileProductComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AgileProductComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
