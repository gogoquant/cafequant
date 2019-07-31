import { SelectuserModule } from './selectuser.module';

describe('SelectuserModule', () => {
  let selectuserModule: SelectuserModule;

  beforeEach(() => {
    selectuserModule = new SelectuserModule();
  });

  it('should create an instance', () => {
    expect(selectuserModule).toBeTruthy();
  });
});
