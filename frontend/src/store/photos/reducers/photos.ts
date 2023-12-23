import { ActionStatus, Photo } from "@core/domain/domain";
import { createSlice } from "@reduxjs/toolkit";
import { State } from "@store/store.ts";
import { getActionLoadPhotoList } from "@store/photos/creators/photo";

const initialState = {
  "photos/list": {
    status: ActionStatus.IDLE,
    data: [],
  } as State<Photo[]>,
};

export type PhotoState = typeof initialState;

export const photoSlice = createSlice({
  name: "photos",
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(getActionLoadPhotoList.pending, (state: PhotoState) => {
      state["photos/list"]["status"] = ActionStatus.LOADING;
    });
    builder.addCase(
      getActionLoadPhotoList.fulfilled,
      (state: PhotoState, action: any) => {
        state["photos/list"] = {
          ...state["photos/list"],
          status: ActionStatus.DONE,
          data: [
            ...(state["photos/list"].data || []),
            ...action.payload.data.items,
          ],
        };
      }
    );
    builder.addCase(getActionLoadPhotoList.rejected, (state: PhotoState) => {
      state["photos/list"]["status"] = ActionStatus.ERROR;
    });
  },
});

export default photoSlice.reducer;
