import { MceModule } from './mce.module';

describe('MceModule', () => {
  let mceModule: MceModule;

  beforeEach(() => {
    mceModule = new MceModule();
  });

  it('should create an instance', () => {
    expect(mceModule).toBeTruthy();
  });
});
