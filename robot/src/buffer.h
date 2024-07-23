#ifndef BUFFER_H
#define BUFFER_H

#include "arm.h"

#define UART_BUFFER_SIZE 32

void clear_buffer(uint8_t *buffer);

RESULT_CODE load_joints_angles_from_buffer(uint8_t *buffer, size_t buffer_len, JointsAngles *joints_angles);

void read_arm_new_speed_from_buffer(uint8_t *buffer, float *speed);

size_t load_result_with_joints_angles_to_buffer(uint8_t *buffer, RESULT_CODE code, JointsAngles *joints_angles);

size_t load_result_code_to_buffer(uint8_t *buffer, RESULT_CODE code);

#endif