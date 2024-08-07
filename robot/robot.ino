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

uint8_t buffer[UART_BUFFER_SIZE] = {0};
size_t loaded_bytes;
RESULT_CODE result_code;

void send_result(size_t size);

void setup() {
  // configure uart
  Serial.begin(BAUD_RATE, SERIAL_8E1);

  // initialize arm motors
  initialize_arm_motors(&arm);
}

void loop() {
  if (Serial.available() > 0) {
    int loaded_bytes = Serial.readBytesUntil(END_OF_TRANSMISSION, buffer, UART_BUFFER_SIZE);
    switch (buffer[0])
    {
      case SET_NEW_ARM_POSITION: {
        JointsAngles *translations = (JointsAngles*) malloc(sizeof(JointsAngles));
        result_code = load_joints_angles_from_buffer(buffer, loaded_bytes, translations);
        clear_buffer(buffer);

        if (result_code != RESULT_OK) {
          loaded_bytes = load_result_code_to_buffer(buffer, result_code);
          send_result(loaded_bytes);
          clear_buffer(buffer);
          free(translations);
          break;
        }

        JointsAngles *fallback = (JointsAngles*) malloc(sizeof(JointsAngles));
        result_code = set_new_arm_position(&arm, translations, fallback);

        if (result_code != RESULT_OK) {
          loaded_bytes = load_result_code_to_buffer(buffer, result_code);
          send_result(loaded_bytes);
          clear_buffer(buffer);

          free(translations);
          free(fallback);
          break;
        }

        loaded_bytes = load_result_with_joints_angles_to_buffer(buffer, result_code, fallback);
        send_result(loaded_bytes);
        clear_buffer(buffer);

        free(translations);
        free(fallback);
        break;
      }
      case SET_ARM_SPEED: {
        float *new_speed = (float *) malloc(sizeof(float));
        read_arm_new_speed_from_buffer(buffer, new_speed);
        clear_buffer(buffer);
        
        result_code = set_arm_speed(&arm, *new_speed);
        loaded_bytes = load_result_code_to_buffer(buffer, result_code);
        send_result(loaded_bytes);
        free(new_speed);
        break;
      }
      case GET_ARM_CURRENT_POSITION: {
        clear_buffer(buffer);
        JointsAngles *current_position = (JointsAngles*) malloc(sizeof(JointsAngles));
        result_code = get_arm_current_position(&arm, current_position);

        if (result_code != RESULT_OK) {
          loaded_bytes = load_result_code_to_buffer(buffer, result_code);
          send_result(loaded_bytes);
          clear_buffer(buffer);
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
        result_code = is_arm_calibrated(&arm);
        loaded_bytes = load_result_code_to_buffer(buffer, result_code);
        send_result(loaded_bytes);
        break;
      }
      case START_CALIBRATION: {
        clear_buffer(buffer);
        if (is_arm_in_move(&arm)) {
          loaded_bytes = load_result_code_to_buffer(buffer, RESULT_ARM_IN_MOVE);
          send_result(loaded_bytes);
          break;
        } 

        set_arm_mode(&arm, ARM_CALIBRATION_MODE);
        set_arm_calibration(&arm, false);

        loaded_bytes = load_result_code_to_buffer(buffer, RESULT_OK);
        send_result(loaded_bytes);
        break;
      }
      case FINISH_CALIBRATION: {
        clear_buffer(buffer);
        result_code = set_arm_current_position_as_reference(&arm);
        
        if (result_code != RESULT_OK) {
          loaded_bytes = load_result_code_to_buffer(buffer, result_code);
          send_result(loaded_bytes);
          break;
        }

        set_arm_mode(&arm, ARM_NORMAL_MODE);
        set_arm_calibration(&arm, true);
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
        set_arm_mode(&arm, ARM_NORMAL_MODE);
        loaded_bytes = load_result_code_to_buffer(buffer, RESULT_OK);
        send_result(loaded_bytes);
        break;
      }
      case CHECK_ARM_IDLE: {
        clear_buffer(buffer);
        if (is_arm_in_move(&arm)) result_code = RESULT_ARM_IN_MOVE;
        else result_code = RESULT_OK;
        
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
  move_arm_steppers(&arm);
}

void send_result(size_t size) {
  buffer[size] = END_OF_TRANSMISSION;
  Serial.write(buffer, size+1);
}