import { TmpModule } from './tmp.module';

describe('TmpModule', () => {
  let tmpModule: TmpModule;

  beforeEach(() => {
    tmpModule = new TmpModule();
  });

  it('should create an instance', () => {
    expect(tmpModule).toBeTruthy();
  });
});
