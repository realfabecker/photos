import { Container as InversifyContainer } from "inversify";
import { Types } from "@core/ports/ports.ts";

import { LocalAuthService } from "@core/adapters/LocalAuthService.ts";
import { PicsumPhotoService } from "@core/adapters/PicsumPhotoService.ts";
import { ProviderEnum } from "@core/domain/domain.ts";
import { SintesePhotoService } from "@core/adapters/SintesePhotoService.ts";
import { SinteseAuthService } from "@core/adapters/SinteseAuthService.ts";

const m = location.search.match(/\?provider=(?<provider>.*)&?/);
const p = m?.groups?.["provider"] || ProviderEnum.Picsum;

export const container = new InversifyContainer();
if (p === ProviderEnum.Lambda) {
  container.bind(Types.PhotoService).to(SintesePhotoService);
  container.bind(Types.AuthService).to(SinteseAuthService);
} else {
  container.bind(Types.PhotoService).to(PicsumPhotoService);
  container.bind(Types.AuthService).to(LocalAuthService);
}
