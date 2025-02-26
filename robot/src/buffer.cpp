#include "buffer.h"

const uint8_t ACTION_ID_OFFSET = 0;
const uint8_t ACTION_ID_SIZE = 1;
const uint8_t X_JOINT_ANGLE_OFFSET = ACTION_ID_OFFSET + ACTION_ID_SIZE;
const uint8_t X_JOINT_ANGLE_SIZE = 4;
const uint8_t Y_JOINT_ANGLE_OFFSET = X_JOINT_ANGLE_OFFSET + X_JOINT_ANGLE_SIZE;
const uint8_t Y_JOINT_ANGLE_SIZE = 4;
const uint8_t Z_JOINT_ANGLE_OFFSET = Y_JOINT_ANGLE_OFFSET + Y_JOINT_ANGLE_SIZE;
const uint8_t Z_JOINT_ANGLE_SIZE = 4;
const uint8_t V_JOINT_ANGLE_OFFSET = Z_JOINT_ANGLE_OFFSET + Z_JOINT_ANGLE_SIZE;
const uint8_t V_JOINT_ANGLE_SIZE = 4;
const uint8_t W_JOINT_ANGLE_OFFSET = V_JOINT_ANGLE_OFFSET + V_JOINT_ANGLE_SIZE;
const uint8_t W_JOINT_ANGLE_SIZE = 4;

const uint8_t RESULT_CODE_SIZE = 1;

const uint8_t TOTAL_ACTION_WITH_JOINTS_ANGLES_SIZE = ACTION_ID_SIZE + X_JOINT_ANGLE_SIZE + Y_JOINT_ANGLE_SIZE + Z_JOINT_ANGLE_SIZE + V_JOINT_ANGLE_SIZE + W_JOINT_ANGLE_SIZE;
const uint8_t TOTAL_RESULT_WITH_JOINTS_ANGLES_SIZE = RESULT_CODE_SIZE + X_JOINT_ANGLE_SIZE + Y_JOINT_ANGLE_SIZE + Z_JOINT_ANGLE_SIZE + V_JOINT_ANGLE_SIZE + W_JOINT_ANGLE_SIZE;

const uint8_t ARM_SPEED_OFFSET = ACTION_ID_OFFSET + ACTION_ID_SIZE;
const uint8_t ARM_SPEED_SIZE = 4;


void clear_buffer(uint8_t *buffer) {
  memset(buffer, 0, UART_BUFFER_SIZE);
}

RESULT_CODE load_joints_angles_from_buffer(
    uint8_t *buffer,
    size_t buffer_len,
    JointsAngles *joints_angles
) {
    if (buffer_len < TOTAL_ACTION_WITH_JOINTS_ANGLES_SIZE) {
        return RESULT_INVALID_NUMBER_OF_PARAMETERS;
    }

    memcpy(&(joints_angles->x), buffer+X_JOINT_ANGLE_OFFSET, X_JOINT_ANGLE_SIZE);
    memcpy(&(joints_angles->y), buffer+Y_JOINT_ANGLE_OFFSET, Y_JOINT_ANGLE_SIZE);
    memcpy(&(joints_angles->z), buffer+Z_JOINT_ANGLE_OFFSET, Z_JOINT_ANGLE_SIZE);
    memcpy(&(joints_angles->v), buffer+V_JOINT_ANGLE_OFFSET, V_JOINT_ANGLE_SIZE);
    memcpy(&(joints_angles->w), buffer+W_JOINT_ANGLE_OFFSET, W_JOINT_ANGLE_SIZE);

    return RESULT_OK;
}

void read_arm_new_speed_from_buffer(uint8_t *buffer, float *speed) {
    memcpy(speed, buffer+ARM_SPEED_OFFSET, ARM_SPEED_SIZE);
}

size_t load_result_with_joints_angles_to_buffer(uint8_t *buffer, RESULT_CODE code, JointsAngles *joints_angles) {
    buffer[0] = code;
    memcpy(buffer+X_JOINT_ANGLE_OFFSET, &(joints_angles->x), X_JOINT_ANGLE_SIZE);
    memcpy(buffer+Y_JOINT_ANGLE_OFFSET, &(joints_angles->y), Y_JOINT_ANGLE_SIZE);
    memcpy(buffer+Z_JOINT_ANGLE_OFFSET, &(joints_angles->z), Z_JOINT_ANGLE_SIZE);
    memcpy(buffer+V_JOINT_ANGLE_OFFSET, &(joints_angles->v), V_JOINT_ANGLE_SIZE);
    memcpy(buffer+W_JOINT_ANGLE_OFFSET, &(joints_angles->w), W_JOINT_ANGLE_SIZE);

    return TOTAL_RESULT_WITH_JOINTS_ANGLES_SIZE;
}

size_t load_result_code_to_buffer(uint8_t *buffer, RESULT_CODE code) {
    buffer[0] = code;

    return RESULT_CODE_SIZE;
}

void add_number_of_loaded_bytes_at_the_buffer_beginning(uint8_t *buffer, size_t number_of_loaded_bytes) {
    memmove(buffer+sizeof(uint8_t), buffer, number_of_loaded_bytes*sizeof(uint8_t));
    buffer[0] = (uint8_t)number_of_loaded_bytes;
}