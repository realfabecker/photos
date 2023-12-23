import { createAsyncThunk } from "@reduxjs/toolkit";
import { Container } from "inversify";
import { IPhotoService, Types } from "@core/ports/ports.ts";

export const getActionLoadPhotoList = createAsyncThunk(
  "transactions/list",
  async (opts: { page: number; limit: number }, { extra }) => {
    const container = (<any>extra).container as Container;
    const service = container.get<IPhotoService>(Types.PhotoService);
    return service.fetchPhotos({ page: opts.page, limit: opts.limit });
  }
);
