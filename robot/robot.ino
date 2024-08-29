#include "src/arm.h"
#include "src/buffer.h"

#define BAUD_RATE 115200
#define END_OF_TRANSMISSION 0x04


AccelStepper y_stepper(AccelStepper::DRIVER, Y_STEP_PIN, Y_DIRECTION_PIN);
AccelStepper z_stepper(AccelStepper::DRIVER, Z_STEP_PIN, Z_DIRECTION_PIN);
AccelStepper x_stepper(AccelStepper::DRIVER, X_STEP_PIN, X_DIRECTION_PIN);
Servo v_servo;
Servo w_servo;

Arm arm = {
  &x_stepper,
  &y_stepper, 
  &z_stepper, 
  &v_servo, 
  &w_servo, 
  ArmState{false, ARM_NORMAL_MODE}
};

uint8_t bytes_to_read = 0;
uint8_t buffer[UART_BUFFER_SIZE] = {0};
size_t loaded_bytes;
RESULT_CODE result_code;

void send_result(size_t size);

void setup() {
  // configure uart
  Serial.begin(BAUD_RATE, SERIAL_8E1);

  // initialize arm motors
  arm.initialize_motors();
}

void loop() {
  if (Serial.available() > 0) {
    Serial.readBytes(&bytes_to_read, 1);

    size_t loaded_bytes = 0;
    while(loaded_bytes < bytes_to_read) {
      size_t n = Serial.readBytes(buffer, bytes_to_read);
      loaded_bytes += n;
    }

    switch (buffer[0])
    {
      case SET_NEW_ARM_POSITION: {
        JointsAngles *translations = (JointsAngles*) malloc(sizeof(JointsAngles));
        result_code = load_joints_angles_from_buffer(buffer, loaded_bytes, translations);
        clear_buffer(buffer);

        if (result_code != RESULT_OK) {
          loaded_bytes = load_result_code_to_buffer(buffer, result_code);
          send_result(loaded_bytes);
          free(translations);
          break;
        }

        JointsAngles *fallback = (JointsAngles*) malloc(sizeof(JointsAngles));
        result_code = arm.set_new_position(translations, fallback);

        if (result_code != RESULT_OK) {
          loaded_bytes = load_result_code_to_buffer(buffer, result_code);
          send_result(loaded_bytes);

          free(translations);
          free(fallback);
          break;
        }

        loaded_bytes = load_result_with_joints_angles_to_buffer(buffer, result_code, fallback);
        send_result(loaded_bytes);

        free(translations);
        free(fallback);
        break;
      }
      case SET_ARM_SPEED: {
        float *new_speed = (float *) malloc(sizeof(float));
        read_arm_new_speed_from_buffer(buffer, new_speed);
        clear_buffer(buffer);
        
        result_code = arm.set_speed(*new_speed);
        loaded_bytes = load_result_code_to_buffer(buffer, result_code);
        send_result(loaded_bytes);
        free(new_speed);
        break;
      }
      case GET_ARM_CURRENT_POSITION: {
        clear_buffer(buffer);
        JointsAngles *current_position = (JointsAngles*) malloc(sizeof(JointsAngles));
        result_code = arm.get_current_position(current_position);

        if (result_code != RESULT_OK) {
          loaded_bytes = load_result_code_to_buffer(buffer, result_code);
          send_result(loaded_bytes);
          free(current_position);
          break;
        }

        loaded_bytes = load_result_with_joints_angles_to_buffer(buffer, result_code, current_position);
        send_result(loaded_bytes);
        free(current_position);
        break;
      }
      case CHECK_ARM_CALIBRATION: { 
        clear_buffer(buffer);
        result_code = arm.is_calibrated();
        loaded_bytes = load_result_code_to_buffer(buffer, result_code);
        send_result(loaded_bytes);
        break;
      }
      case START_CALIBRATION: {
        clear_buffer(buffer);
        if (arm.is_in_move()) {
          loaded_bytes = load_result_code_to_buffer(buffer, RESULT_ARM_IN_MOVE);
          send_result(loaded_bytes);
          break;
        } 

        arm.set_mode(ARM_CALIBRATION_MODE);
        arm.set_calibration(false);

        loaded_bytes = load_result_code_to_buffer(buffer, RESULT_OK);
        send_result(loaded_bytes);
        break;
      }
      case FINISH_CALIBRATION: {
        clear_buffer(buffer);
        result_code = arm.set_current_position_as_reference();
        
        if (result_code != RESULT_OK) {
          loaded_bytes = load_result_code_to_buffer(buffer, result_code);
          send_result(loaded_bytes);
          break;
        }

        arm.set_mode(ARM_NORMAL_MODE);
        arm.set_calibration(true);
        
        loaded_bytes = load_result_code_to_buffer(buffer, result_code);
        send_result(loaded_bytes);
        break;
      }
      case ABORT_CALIBRATION: {
        clear_buffer(buffer);
        if (arm.state.mode != ARM_CALIBRATION_MODE) {
          loaded_bytes = load_result_code_to_buffer(buffer, RESULT_ARM_NOT_IN_CALIBRATION_MODE);
          send_result(loaded_bytes);
          break;
        }
        arm.set_mode(ARM_NORMAL_MODE);
        loaded_bytes = load_result_code_to_buffer(buffer, RESULT_OK);
        send_result(loaded_bytes);
        break;
      }
      case CHECK_ARM_IDLE: {
        clear_buffer(buffer);
        if (arm.is_in_move()) result_code = RESULT_ARM_IN_MOVE;
        else result_code = RESULT_OK;
        
        loaded_bytes = load_result_code_to_buffer(buffer, result_code);
        send_result(loaded_bytes);
        break;
      }
      case OPEN_GRIPPER: {
        clear_buffer(buffer);
        if (arm.is_calibrated()) {
          arm.open_gripper();
          result_code = RESULT_OK;
        }
        else result_code = RESULT_ARM_NOT_CALIBRATED;

        loaded_bytes = load_result_code_to_buffer(buffer, result_code);
        send_result(loaded_bytes);
        break;
      }
      case CLOSE_GRIPPER: {
        clear_buffer(buffer);
        if (arm.is_calibrated()) {
          arm.close_gripper();
          result_code = RESULT_OK;
        }
        else result_code = RESULT_ARM_NOT_CALIBRATED;

        loaded_bytes = load_result_code_to_buffer(buffer, result_code);
        send_result(loaded_bytes);
        break;
      }
      default: {
        loaded_bytes = load_result_code_to_buffer(buffer, RESULT_UNKNOWN_ACTION);
        send_result(loaded_bytes);
        break;
      }
    }
    clear_buffer(buffer);
  }
  arm.move_steppers();
}

void send_result(size_t size) {
  add_number_of_loaded_bytes_at_the_buffer_beginning(buffer, size);
  Serial.write(buffer, size+1);
}