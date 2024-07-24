import tkinter as tk
import cv2
from PIL import Image, ImageTk
from tkinter import ttk
from websocket_client import WsClient


class RobotGUI:
    def __init__(self, root: tk.Tk):
        self.gripper_button = None
        self.send_xyz_button = None
        self.send_joints_button = None
        self.right_sliders = None
        self.right_entries = None
        self.left_sliders = None
        self.left_entries = None

        self.root = root
        self.root.title("Wariat Robot - GUI")

        self.gripper_status = self.get_initial_gripper_status()

        self.left_frame = ttk.Frame(self.root)
        self.left_frame.grid(row=0, column=0, padx=10, pady=10, sticky="n")

        self.video_frame = ttk.Frame(self.root)
        self.video_frame.grid(row=0, column=1, padx=10, pady=10)
        self.create_left_panel()

        self.right_frame = ttk.Frame(self.root)
        self.right_frame.grid(row=0, column=2, padx=10, pady=10, sticky="n")
        self.create_right_panel()

        self.video_label = ttk.Label(self.video_frame)
        self.video_label.grid(row=0, column=0)
        self.cap = cv2.VideoCapture(0)  # można tu dać stream RTSP
        self.update_video_frame()

        self.websocket_client = WsClient("ws://localhost:8765")

    def get_initial_gripper_status(self) -> bool:
        return True

    def create_left_panel(self) -> None:
        self.left_sliders = []
        self.left_entries = []
        slider_ranges = [
            (0, 360),  # Range for Joint 0
            (-95, 110),  # Range for Joint 1
            (-90, 130),  # Range for Joint 2
            (0, 360),  # Range for Joint 3
            (-100, 100)  # Range for Joint 4
        ]

        for i, (min_val, max_val) in enumerate(slider_ranges):
            frame = ttk.Frame(self.left_frame)
            frame.grid(row=i, column=0, padx=5, pady=5)

            slider = tk.Scale(
                frame,
                from_=min_val,
                to=max_val,
                orient="horizontal",
                length=300,
                label=f"Joint {i}",
                command=lambda value, i=i: self.update_entry_from_slider(value, i, 'left')
            )
            slider.grid(row=0, column=0, padx=5, pady=5)
            self.left_sliders.append(slider)

            entry = ttk.Entry(frame)
            entry.grid(row=0, column=1, padx=5, pady=5)
            entry.bind('<Return>', lambda event, i=i: self.update_slider_from_entry(event, i, 'left'))
            self.left_entries.append(entry)

        self.send_joints_button = ttk.Button(
            self.left_frame, text="Send commands", command=self.send_joints_commands
        )
        self.send_joints_button.grid(row=6, column=0, padx=5, pady=5)

    def create_right_panel(self) -> None:
        self.right_sliders = []
        self.right_entries = []
        slider_ranges = [
            (0, 100),  # Range for Slider 1
            (0, 100),  # Range for Slider 2
            (0, 100)  # Range for Slider 3
        ]

        for i, (min_val, max_val) in enumerate(slider_ranges):
            frame = ttk.Frame(self.right_frame)
            frame.grid(row=i, column=0, padx=5, pady=5)

            slider = tk.Scale(
                frame,
                from_=min_val,
                to=max_val,
                orient="horizontal",
                length=300,
                label=f"Slider {i + 1}",
                command=lambda value, i=i: self.update_entry_from_slider(value, i, 'right')
            )
            slider.grid(row=0, column=0, padx=5, pady=5)
            self.right_sliders.append(slider)

            entry = ttk.Entry(frame)
            entry.grid(row=0, column=1, padx=5, pady=5)
            entry.bind('<Return>', lambda event, i=i: self.update_slider_from_entry(event, i, 'right'))
            self.right_entries.append(entry)

        self.gripper_button = ttk.Button(
            self.right_frame, text=self.get_gripper_text(), command=self.toggle_gripper
        )
        self.gripper_button.grid(row=3, column=0, padx=5, pady=5)

        self.send_xyz_button = ttk.Button(
            self.right_frame, text="Send commands", command=self.send_xyz_commands
        )
        self.send_xyz_button.grid(row=4, column=0, padx=5, pady=5)

    def update_entry_from_slider(self, value, index, side):
        if side == 'left':
            self.left_entries[index].delete(0, tk.END)
            self.left_entries[index].insert(0, str(value))
        else:
            self.right_entries[index].delete(0, tk.END)
            self.right_entries[index].insert(0, str(value))

    def update_slider_from_entry(self, event, index, side):
        if side == 'left':
            value = self.left_entries[index].get()
            self.left_sliders[index].set(value)
        else:
            value = self.right_entries[index].get()
            self.right_sliders[index].set(value)

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
