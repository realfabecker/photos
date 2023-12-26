import { ActionStatus, Photo } from "@core/domain/domain";
import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { State } from "@store/store.ts";
import {
  getActionLoadPhotoList,
  getActionPhotosCreate,
  getActionUploadGetUrl,
  getActionUploadRequest,
} from "@store/photos/creators/photo";

const initialState = {
  "photos/list": {
    status: ActionStatus.IDLE,
    data: [],
    has_more: false,
    page_token: "",
  } as State<Photo[]>,
  "photos/add": {
    status: ActionStatus.IDLE,
    modal: "closed",
  } as State,
  "photos/upload": {
    status: ActionStatus.IDLE,
    upload_url: "",
  } as State,
};

export type PhotoState = typeof initialState;

export const photoSlice = createSlice({
  name: "photos",
  initialState: initialState,
  reducers: {
    modal_set: (state, action: PayloadAction<{ state: string }>) => {
      state["photos/add"].modal = action.payload.state as "open" | "closed";
    },
  },
  extraReducers: (builder) => {
    /**
     * Photos List
     */
    builder.addCase(getActionLoadPhotoList.pending, (state: PhotoState) => {
      state["photos/list"]["status"] = ActionStatus.LOADING;
    });
    builder.addCase(
      getActionLoadPhotoList.fulfilled,
      (state: PhotoState, action: any) => {
        if (action.meta.arg.page === 1) {
          state["photos/list"] = {
            ...state["photos/list"],
            status: ActionStatus.DONE,
            data: action.payload.data.items,
            page_token: action.payload.data?.page_token,
            has_more: action.payload.data.has_more,
          };
        } else {
          state["photos/list"] = {
            ...state["photos/list"],
            status: ActionStatus.DONE,
            data: [
              ...(state["photos/list"].data || []),
              ...action.payload.data.items,
            ],
            page_token: action.payload.data?.page_token,
            has_more: action.payload.data.has_more,
          };
        }
      }
    );
    builder.addCase(getActionLoadPhotoList.rejected, (state: PhotoState) => {
      state["photos/list"]["status"] = ActionStatus.ERROR;
    });

    /**
     * Photos Get Upload Url
     */
    builder.addCase(getActionUploadGetUrl.pending, (state: PhotoState) => {
      state["photos/upload"]["status"] = ActionStatus.LOADING;
    });
    builder.addCase(
      getActionUploadGetUrl.fulfilled,
      (state: PhotoState, action: any) => {
        state["photos/upload"] = {
          ...state["photos/upload"],
          status: ActionStatus.DONE,
          upload_url: action.payload,
        };
      }
    );
    builder.addCase(getActionUploadGetUrl.rejected, (state: PhotoState) => {
      state["photos/upload"]["status"] = ActionStatus.ERROR;
      state["photos/upload"]["error"] = {
        message: "Não foi possível carregar a imagem",
      };
    });

    /**
     * Photos Upload Request
     */
    builder.addCase(getActionUploadRequest.pending, (state: PhotoState) => {
      state["photos/upload"] = {
        ...state["photos/upload"],
        status: ActionStatus.LOADING,
      };
    });
    builder.addCase(getActionUploadRequest.fulfilled, (state: PhotoState) => {
      state["photos/upload"] = {
        ...state["photos/upload"],
        status: ActionStatus.DONE,
      };
    });
    builder.addCase(getActionUploadRequest.rejected, (state: PhotoState) => {
      state["photos/upload"] = {
        status: ActionStatus.ERROR,
        error: {
          message: "Erro ao carregar a imagem",
        },
      };
    });

    /**
     * Photos Add
     */
    builder.addCase(getActionPhotosCreate.pending, (state: PhotoState) => {
      state["photos/add"]["status"] = ActionStatus.LOADING;
    });
    builder.addCase(getActionPhotosCreate.fulfilled, (state: PhotoState) => {
      state["photos/add"] = {
        status: ActionStatus.DONE,
      };
    });
    builder.addCase(getActionPhotosCreate.rejected, (state: PhotoState) => {
      state["photos/add"]["status"] = ActionStatus.ERROR;
    });
  },
});

export default photoSlice.reducer;
