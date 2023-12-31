import "reflect-metadata";
import { injectable } from "inversify";
import { IPhotoService } from "@core/ports/ports.ts";
import { PagedDTO, Photo, ResponseDTO } from "@core/domain/domain.ts";

type PicsumPhoto = {
  id: string;
  author: string;
  download_url: string;
  height: string;
  width: string;
};

@injectable()
export class PicsumPhotoService implements IPhotoService {
  async fetchPhotos(opts: {
    page: number;
    limit: number;
    token?: string;
  }): Promise<ResponseDTO<PagedDTO<Photo>>> {
    const res = await fetch(
      `https://picsum.photos/v2/list?page=${opts.page}&limit=${opts.limit}`
    );
    const data = (await res.json()) as PicsumPhoto[];

    const items: Photo[] = data.map((p) => ({
      url: `https://picsum.photos/id/${p.id}/400/225`,
      id: p.id,
      tags: ["tag"],
      title: p.author,
      createdAt: new Date().toISOString(),
    }));

    return {
      status: "success",
      data: {
        items: items,
        page_count: items.length,
        has_more: true,
      },
    };
  }

  async createPhoto(photo: Partial<Photo>): Promise<ResponseDTO<Photo>> {
    console.log({ photo });
    return Promise.reject("not implemented yet");
  }

  async getUploadUrl(name: string): Promise<string> {
    console.log({ name });
    return Promise.reject("not implemented yet");
  }

  async uploadFile(file: File, url: string): Promise<void> {
    console.log({ file, url });
    return Promise.reject("not implemented yet");
  }
}
