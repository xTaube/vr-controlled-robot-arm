import os
import tkinter as tk
import cv2
from PIL import Image, ImageTk
from tkinter import ttk
from websocket_client import WsClient
from dotenv import load_dotenv

load_dotenv()
WEBSOCKET_URL = os.getenv("WEBSOCKET_URL")


class RobotGUI:
    def __init__(self, root: tk.Tk):
        self.set_all_to_0_button = None
        self.gripper_button = None
        self.send_xyz_button = None
        self.send_joints_button = None
        self.right_sliders = None
        self.left_sliders = None
        self.calibration_sliders = None
        self.calibration_button = None
        self.toggle_calibration_button = None
        self.is_calibrating = False

        self.root = root
        self.root.title("Wariat Robot - GUI")

        style = ttk.Style()
        style.configure("Green.TButton", background="green")
        style.configure("Red.TButton", background="red")

        self.gripper_status = self.get_initial_gripper_status()

        self.left_frame = ttk.Frame(self.root)
        self.left_frame.grid(row=0, column=0, padx=10, pady=10, sticky="n")

        self.video_frame = ttk.Frame(self.root)
        self.video_frame.grid(row=0, column=1, padx=10, pady=10)
        self.create_left_panel()

        self.right_frame = ttk.Frame(self.root)
        self.right_frame.grid(row=0, column=2, padx=10, pady=10, sticky="n")
        self.create_right_panel()

        self.calibration_frame = ttk.Frame(self.root)
        self.calibration_frame.grid(row=1, column=1, padx=10, pady=10)

        self.video_label = ttk.Label(self.video_frame)
        self.video_label.grid(row=0, column=0)
        self.cap = cv2.VideoCapture(0)  # można tu dać stream RTSP
        self.update_video_frame()

        self.websocket_client = WsClient(WEBSOCKET_URL)

    def get_initial_gripper_status(self) -> bool:
        return True

    def create_left_panel(self) -> None:
        self.left_sliders = []
        slider_ranges = [
            (-360, 360, "z"),  # Range for Joint z
            (-65, 95, "y"),  # Range for Joint y
            (-90, 95, "x"),  # Range for Joint x
            (0, 180, "v"),  # Range for Joint v
            (5, 180, "w"),  # Range for Joint w
        ]

        for i, (min_val, max_val, axis_name) in enumerate(slider_ranges):
            frame = ttk.Frame(self.left_frame)
            frame.grid(row=i, column=0, padx=5, pady=5)

            min_label = ttk.Label(frame, text=f"{min_val}")
            min_label.grid(row=0, column=0, padx=5, pady=5)

            slider = tk.Scale(
                frame,
                from_=min_val,
                to=max_val,
                orient="horizontal",
                length=300,
                label=f"Joint {axis_name}",
            )
            slider.grid(row=0, column=1, padx=5, pady=5)
            self.left_sliders.append(slider)

            max_label = ttk.Label(frame, text=f"{max_val}")
            max_label.grid(row=0, column=2, padx=5, pady=5)

        self.send_joints_button = ttk.Button(
            self.left_frame, text="Send commands", command=self.send_joints_commands
        )
        self.send_joints_button.grid(row=6, column=0, padx=5, pady=5)

        self.toggle_calibration_button = ttk.Button(
            self.left_frame,
            text="Start Calibrating",
            command=self.toggle_calibration,
            style="TButton",
        )
        self.toggle_calibration_button.grid(row=7, column=0, padx=5, pady=5)

        self.set_all_to_0_button = ttk.Button(
            self.left_frame,
            text="Set all joints to 0",
            command=self.set_all_to_0
        )
        self.set_all_to_0_button.grid(row=8, column=0, padx=5, pady=5)

    def set_all_to_0(self) -> None:
        for slider in self.left_sliders:
            slider.set(0)

    def create_right_panel(self) -> None:
        self.right_sliders = []
        slider_ranges = [
            (0, 100),  # Range for Slider 1
            (0, 100),  # Range for Slider 2
            (0, 100),  # Range for Slider 3
        ]

        for i, (min_val, max_val) in enumerate(slider_ranges):
            frame = ttk.Frame(self.right_frame)
            frame.grid(row=i, column=0, padx=5, pady=5)

            min_label = ttk.Label(frame, text=f"{min_val}")
            min_label.grid(row=0, column=0, padx=5, pady=5)

            slider = tk.Scale(
                frame,
                from_=min_val,
                to=max_val,
                orient="horizontal",
                length=300,
                label=f"Slider {i + 1}",
            )
            slider.grid(row=0, column=1, padx=5, pady=5)
            self.right_sliders.append(slider)

            max_label = ttk.Label(frame, text=f"{max_val}")
            max_label.grid(row=0, column=2, padx=5, pady=5)

        self.gripper_button = ttk.Button(
            self.right_frame, text=self.get_gripper_text(), command=self.toggle_gripper
        )
        self.gripper_button.grid(row=3, column=0, padx=5, pady=5)

        self.send_xyz_button = ttk.Button(
            self.right_frame, text="Send commands", command=self.send_xyz_commands
        )
        self.send_xyz_button.grid(row=4, column=0, padx=5, pady=5)

    def toggle_calibration(self) -> None:
        self.is_calibrating = not self.is_calibrating

        if self.is_calibrating:
            self.websocket_client.send_message("6")
            self.send_joints_button.configure(text="Send calibration")
            self.toggle_calibration_button.configure(
                text="Stop Calibrating", style="Red.TButton"
            )
        else:
            self.websocket_client.send_message("1")
            self.send_joints_button.configure(text="Send commands")
            self.toggle_calibration_button.configure(
                text="Start Calibrating", style="Green.TButton"
            )
            self.set_all_to_0()

    def get_gripper_text(self) -> str:
        return "Open Gripper" if not self.gripper_status else "Close Gripper"

    def toggle_gripper(self) -> None:
        self.gripper_status = not self.gripper_status
        self.gripper_button.configure(text=self.get_gripper_text())

    def send_xyz_commands(self) -> None:
        """To be implemented"""
        pass

    def send_joints_commands(self) -> None:
        command = "3$" + "$".join(str(slider.get()) for slider in self.left_sliders)
        self.websocket_client.send_message(command)

    def update_video_frame(self) -> None:
        ret, frame = self.cap.read()
        if ret:
            width = 600
            height = 400
            dim = (width, height)

            resized_frame = cv2.resize(frame, dim, interpolation=cv2.INTER_AREA)

            frame = cv2.cvtColor(resized_frame, cv2.COLOR_BGR2RGB)
            img = Image.fromarray(frame)
            imgtk = ImageTk.PhotoImage(image=img)
            self.video_label.imgtk = imgtk
            self.video_label.configure(image=imgtk)

        self.root.after(10, self.update_video_frame)

    def __del__(self):
        if self.cap.isOpened():
            self.cap.release()
        self.websocket_client.close_connection()


def main():
    root = tk.Tk()
    RobotGUI(root)
    root.mainloop()


if __name__ == "__main__":
    main()
