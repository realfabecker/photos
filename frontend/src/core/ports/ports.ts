import { PagedDTO, Photo, ResponseDTO } from "../domain/domain.ts";

export const Types = {
  PhotoService: Symbol.for("PhotoService"),
  AuthService: Symbol.for("AuthService"),
};

export interface IPhotoService {
  fetchPhotos(opts: {
    page: number;
    limit: number;
    token?: string;
  }): Promise<ResponseDTO<PagedDTO<Photo>>>;
}

export interface IAuthService {
  login: ({
    email,
    password,
  }: {
    email: string;
    password: string;
  }) => Promise<void>;
  isLoggedIn(): boolean;
  getAccessToken(): string | undefined;
  logout(): void;
}
