import tkinter as tk
from tkinter import ttk
import cv2
from PIL import Image, ImageTk


class RobotGUI:
    def __init__(self, root: tk.Tk):

        self.root = root
        self.root.title("Wariat Robot - GUI")

        self.gripper_button = None
        self.send_xyz = None
        self.send_joints = None
        self.right_sliders = None
        self.left_sliders = None

        self.gripper_status = self.get_initial_gripper_status()

        self.left_frame = ttk.Frame(root)
        self.left_frame.grid(row=0, column=0, padx=10, pady=10, sticky="n")

        self.video_frame = ttk.Frame(root)
        self.video_frame.grid(row=0, column=1, padx=10, pady=10)
        self.create_left_panel()

        self.right_frame = ttk.Frame(root)
        self.right_frame.grid(row=0, column=2, padx=10, pady=10, sticky="n")
        self.create_right_panel()

        # vid capture
        self.video_label = ttk.Label(self.video_frame)
        self.video_label.grid(row=0, column=0)
        self.cap = cv2.VideoCapture(0)  # mozna tu dac stream RTSP
        self.update_video_frame()


    def get_initial_gripper_status(self) -> bool:
        return True

    def create_left_panel(self) -> None:
        self.left_sliders = []
        for i in range(6):
            slider = tk.Scale(
                self.left_frame,
                from_=0,
                to=100,
                orient="horizontal",
                label=f"Joint {i}",
            )
            slider.grid(row=i, column=0, padx=5, pady=5)
            self.left_sliders.append(slider)

        self.send_joints = ttk.Button(
            self.left_frame, text="Send commands", command=self.send_joints_commands
        )
        self.send_joints.grid(row=6, column=0, padx=5, pady=5)
    def create_right_panel(self) -> None:
        self.right_sliders = []
        for i in range(3):
            slider = tk.Scale(
                self.right_frame,
                from_=0,
                to=100,
                orient="horizontal",
                label=f"Slider {i + 1}",
            )
            slider.grid(row=i, column=0, padx=5, pady=5)
            self.right_sliders.append(slider)

        self.gripper_button = ttk.Button(
            self.right_frame, text=self.get_gripper_text(), command=self.toggle_gripper
        )
        self.gripper_button.grid(row=3, column=0, padx=5, pady=5)

        self.send_xyz = ttk.Button(
            self.right_frame, text="Send commands", command=self.send_xyz_commands
        )
        self.send_xyz.grid(row=4, column=0, padx=5, pady=5)

    def get_gripper_text(self) -> str:
        return "Open Gripper" if not self.gripper_status else "Close Gripper"

    def toggle_gripper(self):
        self.gripper_status = not self.gripper_status
        self.gripper_button.configure(text=self.get_gripper_text())

    def send_xyz_commands(self) -> None:
        pass

    def send_joints_commands(self) -> None:
        pass

    def update_video_frame(self) -> None:
        ret, frame = self.cap.read()
        if ret:
            frame = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
            img = Image.fromarray(frame)
            imgtk = ImageTk.PhotoImage(image=img)
            self.video_label.imgtk = imgtk
            self.video_label.configure(image=imgtk)

        self.root.after(10, self.update_video_frame)

    def __del__(self):
        if self.cap.isOpened():
            self.cap.release()


def main():
    root = tk.Tk()
    app = RobotGUI(root)
    root.mainloop()


if __name__ == "__main__":
    main()