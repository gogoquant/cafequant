import { TtaModule } from './tta.module';

describe('TtaModule', () => {
  let ttaModule: TtaModule;

  beforeEach(() => {
    ttaModule = new TtaModule();
  });

  it('should create an instance', () => {
    expect(ttaModule).toBeTruthy();
  });
});
