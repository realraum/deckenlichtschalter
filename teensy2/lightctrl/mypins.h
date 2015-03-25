#ifndef PINS_H
#define PINS_H
#include <avr/io.h>


#define PIN_HIGH(PORT, PIN) PORT |= (1 << PIN)
#define PIN_LOW(PORT, PIN) PORT &= ~(1 << PIN)
#define PINMODE_OUTPUT PIN_HIGH  //just use DDR instead of PORT
#define PINMODE_INPUT PIN_LOW  //just use DDR instead of PORT

#define OP_SETBIT |=
#define OP_CLEARBIT &= ~
#define OP_CHECK &
#define PIN_SW(PORTDDRREG, PIN, OP) PORTDDRREG OP (1 << PIN)

#define HIGHv OP_SETBIT
#define LOWv OP_CLEARBIT

#define RF_DATA_OUT_PIN   PINF6
#define RF_DATA_OUT_PORT  PORTF
#define RF_DATA_OUT_DDR   DDRF

#define NUM_BUTTONS 15

#define BTN_L1_ON_PIN PINB7
#define BTN_L1_ON_PORT PORTB
#define BTN_L1_ON_DDR DDRB

#define BTN_L1_OFF_PIN PIND0
#define BTN_L1_OFF_PORT PORTD
#define BTN_L1_OFF_DDR DDRD

#define BTN_L2_ON_PIN PINB2
#define BTN_L2_ON_PORT PORTB
#define BTN_L2_ON_DDR DDRB

#define BTN_L2_OFF_PIN PINB3
#define BTN_L2_OFF_PORT PORTB
#define BTN_L2_OFF_DDR DDRB

#define BTN_L3_ON_PIN PINB0
#define BTN_L3_ON_PORT PORTB
#define BTN_L3_ON_DDR DDRB

#define BTN_L3_OFF_PIN PINB1
#define BTN_L3_OFF_PORT PORTB
#define BTN_L3_OFF_DDR DDRB

#define BTN_L4_ON_PIN PINC7
#define BTN_L4_ON_PORT PORTC
#define BTN_L4_ON_DDR DDRC

#define BTN_L4_OFF_PIN PIND6
#define BTN_L4_OFF_PORT PORTD
#define BTN_L4_OFF_DDR DDRD

#define BTN_L5_ON_PIN PIND3
#define BTN_L5_ON_PORT PORTD
#define BTN_L5_ON_DDR DDRD

#define BTN_L5_OFF_PIN PINC6
#define BTN_L5_OFF_PORT PORTC
#define BTN_L5_OFF_DDR DDRC

#define BTN_L6_ON_PIN PIND1
#define BTN_L6_ON_PORT PORTD
#define BTN_L6_ON_DDR DDRD

#define BTN_L6_OFF_PIN PIND2
#define BTN_L6_OFF_PORT PORTD
#define BTN_L6_OFF_DDR DDRD

#define BTN_C1_PIN PIND7
#define BTN_C1_PORT PORTD
#define BTN_C1_DDR DDRD

#define BTN_C2_PIN PINB4
#define BTN_C2_PORT PORTB
#define BTN_C2_DDR DDRB

#define BTN_C3_PIN PINB5
#define BTN_C3_PORT PORTB
#define BTN_C3_DDR DDRB

#define RELAIS_PORT PORTF
#define RELAIS_MASK 0x3F
#define RELAIS_DDR DDRF
#define RELAIS_INVERT_MASK 0x03

uint8_t *btns_on_pins_;
volatile uint8_t **btns_on_pinreg_;
uint8_t *btns_off_pins_;
volatile uint8_t **btns_off_pinreg_;
uint8_t *btns_sig_pins_;
volatile uint8_t **btns_sig_pinreg_;

#define TEST_BTN_ON(x) ((*(btns_on_pinreg_[x]) & (1<<btns_on_pins_[x])) == 0)
#define TEST_BTN_OFF(x) ((*(btns_off_pinreg_[x]) & (1<<btns_off_pins_[x])) == 0)
#define TEST_BTN_SIG(x) ((*(btns_sig_pinreg_[x]) & (1<<btns_sig_pins_[x])) == 0)

#define PULLUP_DDRB ((1<<PINB7) | (1<<PINB2) | (1<<PINB3) | (1<<PINB0) | (1<<PINB1) | (1<<PINB4) | (1<<PINB5))
#define PULLUP_DDRC ((1<<PINC7) | (1<<PINC6))
#define PULLUP_DDRD ((1<<PIND0) | (1<<PIND6) | (1<<PIND3) | (1<<PIND1) | (1<<PIND2) | (1<<PIND7))
#define PULLUP_DDRF 0

#endif