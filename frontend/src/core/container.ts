import { Container as InversifyContainer } from "inversify";
import { Types } from "@core/ports/ports.ts";
import { PicsumPhotoService } from "@core/adapters/PicsumPhotoService.ts";
import { LocalAuthService } from "@core/adapters/LocalAuthService.ts";

export const container = new InversifyContainer();
container.bind(Types.PhotoService).to(PicsumPhotoService);
container.bind(Types.AuthService).to(LocalAuthService);
