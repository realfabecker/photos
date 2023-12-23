import { PagedDTO, Photo, ResponseDTO } from "../domain/domain.ts";

export const Types = {
  PhotoService: Symbol.for("PhotoService"),
};

export interface IPhotoService {
  fetchPhotos(opts: {
    page: number;
    limit: number;
  }): Promise<ResponseDTO<PagedDTO<Photo>>>;
}
