import { CommonsModule } from './commons.module';

describe('CommonModule', () => {
  let commonModule: CommonsModule;

  beforeEach(() => {
    commonModule = new CommonsModule();
  });

  it('should create an instance', () => {
    expect(commonModule).toBeTruthy();
  });
});
