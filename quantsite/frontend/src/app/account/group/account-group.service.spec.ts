import { TestBed, inject } from '@angular/core/testing';

import { AccountGroupService } from './account-group.service';

describe('AccountGroupService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [AccountGroupService]
    });
  });

  it('should be created', inject([AccountGroupService], (service: AccountGroupService) => {
    expect(service).toBeTruthy();
  }));
});
