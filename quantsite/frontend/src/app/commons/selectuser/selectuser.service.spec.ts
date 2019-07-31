import { TestBed, inject } from '@angular/core/testing';

import { SelectuserService } from './selectuser.service';

describe('SelectuserService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [SelectuserService]
    });
  });

  it('should be created', inject([SelectuserService], (service: SelectuserService) => {
    expect(service).toBeTruthy();
  }));
});
