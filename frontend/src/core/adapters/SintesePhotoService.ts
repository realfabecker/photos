import "reflect-metadata";
import { inject, injectable } from "inversify";
import { type IAuthService, IPhotoService, Types } from "@core/ports/ports.ts";
import { PagedDTO, Photo, ResponseDTO } from "@core/domain/domain.ts";

type SintesePaged<T> = {
  status: "success" | "error";
  data: {
    page_count: number;
    items: T[];
    has_more: boolean;
  };
};

type SintesePhoto = {
  photoId: string;
  userId: string;
  fileName: string;
  title: string;
  url: string;
  createdAt: string;
};

@injectable()
export class SintesePhotoService implements IPhotoService {
  constructor(
    @inject(Types.AuthService) private readonly auth: IAuthService,
    private readonly baseUrl: string = import.meta.env.VITE_API_BASE_URL
  ) {}
  async fetchPhotos(opts: {
    page: number;
    limit: number;
  }): Promise<ResponseDTO<PagedDTO<Photo>>> {
    console.log(opts);
    const res = await fetch(`${this.baseUrl}/midia?limit=10&created_at=2023`, {
      headers: { Authorization: `Bearer ${this.auth.getAccessToken()}` },
    });
    const data = (await res.json()) as SintesePaged<SintesePhoto>;
    const items: Photo[] = data.data.items.map((p) => ({
      url: p.url,
      id: p.photoId,
      tags: ["tag"],
      title: p.title,
      createdAt: p.createdAt,
    }));
    return {
      status: "success",
      data: {
        items: items,
        page_count: data.data.page_count,
        has_more: data.data.has_more,
      },
    };
  }
}
