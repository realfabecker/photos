import { Container as InversifyContainer } from "inversify";
import { Types } from "@core/ports/ports.ts";
import { PhotoPicsumService } from "@core/adapters/PhotoPicsumService.ts";

export const container = new InversifyContainer();
container.bind(Types.PhotoService).to(PhotoPicsumService);
