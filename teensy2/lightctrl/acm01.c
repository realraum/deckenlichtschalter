/*
 *  spreadspace avr utils
 *
 *
 *  Copyright (C) 2013-2014 Christian Pointner <equinox@spreadspace.org>
 *
 *  This file is part of spreadspace avr utils.
 *
 *  spreadspace avr utils is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  any later version.
 *
 *  spreadspace avr utils is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with spreadspace avr utils. If not, see <http://www.gnu.org/licenses/>.
 */


#include <avr/io.h>
#include <avr/wdt.h>
#include <avr/interrupt.h>
#include <avr/power.h>

#include "util.h"
/*
             LUFA Library
     Copyright (C) Dean Camera, 2012.

  dean [at] fourwalledcubicle [dot] com
           www.lufa-lib.org
*/
#include <LUFA/Drivers/USB/USB.h>
#include <LUFA/Drivers/Misc/RingBuffer.h>
#include <LUFA/Drivers/Peripheral/Serial.h>
#include "lufa-descriptor-usbdualserial.h"

FILE usb1_stream_;
FILE usb2_stream_;

USB_ClassInfo_CDC_Device_t VirtualSerial1_CDC_Interface =
  {
    .Config =
      {
        .ControlInterfaceNumber           = 0,

        .DataINEndpointNumber             = CDC1_TX_EPNUM,
        .DataINEndpointSize               = CDC_TXRX_EPSIZE,
        .DataINEndpointDoubleBank         = false,

        .DataOUTEndpointNumber            = CDC1_RX_EPNUM,
        .DataOUTEndpointSize              = CDC_TXRX_EPSIZE,
        .DataOUTEndpointDoubleBank        = false,

        .NotificationEndpointNumber       = CDC1_NOTIFICATION_EPNUM,
        .NotificationEndpointSize         = CDC_NOTIFICATION_EPSIZE,
        .NotificationEndpointDoubleBank   = false,
      },
  };

USB_ClassInfo_CDC_Device_t VirtualSerial2_CDC_Interface =
  {
    .Config =
      {
        .ControlInterfaceNumber           = 2,

        .DataINEndpointNumber             = CDC2_TX_EPNUM,
        .DataINEndpointSize               = CDC_TXRX_EPSIZE,
        .DataINEndpointDoubleBank         = false,

        .DataOUTEndpointNumber            = CDC2_RX_EPNUM,
        .DataOUTEndpointSize              = CDC_TXRX_EPSIZE,
        .DataOUTEndpointDoubleBank        = false,

        .NotificationEndpointNumber       = CDC2_NOTIFICATION_EPNUM,
        .NotificationEndpointSize         = CDC_NOTIFICATION_EPSIZE,
        .NotificationEndpointDoubleBank   = false,
      },
  };

static RingBuffer_t USBtoUSART_Buffer;
static uint8_t      USBtoUSART_Buffer_Data[128];
static RingBuffer_t USARTtoUSB_Buffer;
static uint8_t      USARTtoUSB_Buffer_Data[128];


void EVENT_USB_Device_ConfigurationChanged(void)
{
  CDC_Device_ConfigureEndpoints(&VirtualSerial1_CDC_Interface);
  CDC_Device_ConfigureEndpoints(&VirtualSerial2_CDC_Interface);
}

void EVENT_USB_Device_ControlRequest(void)
{
  CDC_Device_ProcessControlRequest(&VirtualSerial1_CDC_Interface);
  CDC_Device_ProcessControlRequest(&VirtualSerial2_CDC_Interface);
}

void EVENT_CDC_Device_LineEncodingChanged(USB_ClassInfo_CDC_Device_t* const CDCInterfaceInfo)
{
  if(CDCInterfaceInfo != &VirtualSerial1_CDC_Interface)
    return;

  uint8_t ConfigMask = 0;

  switch (CDCInterfaceInfo->State.LineEncoding.ParityType)
  {
    case CDC_PARITY_Odd:
      ConfigMask = ((1 << UPM11) | (1 << UPM10));
      break;
    case CDC_PARITY_Even:
      ConfigMask = (1 << UPM11);
      break;
  }

  if (CDCInterfaceInfo->State.LineEncoding.CharFormat == CDC_LINEENCODING_TwoStopBits)
    ConfigMask |= (1 << USBS1);

  switch (CDCInterfaceInfo->State.LineEncoding.DataBits)
  {
    case 6:
      ConfigMask |= (1 << UCSZ10);
      break;
    case 7:
      ConfigMask |= (1 << UCSZ11);
      break;
    case 8:
      ConfigMask |= ((1 << UCSZ11) | (1 << UCSZ10));
      break;
  }

  /* Must turn off USART before reconfiguring it, otherwise incorrect operation may occur */
  UCSR1B = 0;
  UCSR1A = 0;
  UCSR1C = 0;

  /* Set the new baud rate before configuring the USART */
  UBRR1  = SERIAL_2X_UBBRVAL(CDCInterfaceInfo->State.LineEncoding.BaudRateBPS);

  /* Reconfigure the USART in double speed mode for a wider baud rate range at the expense of accuracy */
  UCSR1C = ConfigMask;
  UCSR1A = (1 << U2X1);
  //UCSR1B = ((1 << RXCIE1) | (1 << TXEN1) | (1 << RXEN1));
  //do not enable RXEN1 and RXCIE1, to free pin PIND2
  UCSR1B = ((1 << RXCIE1) | (1 << TXEN1));
}


ISR(USART1_RX_vect, ISR_BLOCK)
{
  uint8_t ReceivedByte = UDR1;

  if (USB_DeviceState == DEVICE_STATE_Configured)
    RingBuffer_Insert(&USARTtoUSB_Buffer, ReceivedByte);
}


void bzero (uint8_t *to, int count)
{
  while (count-- > 0)
    {
      *to++ = 0;
    }
}

void usbserial_init(void)
{
  bzero((uint8_t*)&USBtoUSART_Buffer_Data,sizeof(USBtoUSART_Buffer_Data));
  bzero((uint8_t*)&USARTtoUSB_Buffer_Data,sizeof(USARTtoUSB_Buffer_Data));
  RingBuffer_InitBuffer(&USBtoUSART_Buffer, USBtoUSART_Buffer_Data, sizeof(USBtoUSART_Buffer_Data));
  RingBuffer_InitBuffer(&USARTtoUSB_Buffer, USARTtoUSB_Buffer_Data, sizeof(USARTtoUSB_Buffer_Data));
  TCCR0B = (1 << CS02);
}

void usbserial_task(void)
{
    if (!(RingBuffer_IsFull(&USBtoUSART_Buffer))) {
      int16_t ReceivedByte = CDC_Device_ReceiveByte(&VirtualSerial1_CDC_Interface);
      if (!(ReceivedByte < 0))
        RingBuffer_Insert(&USBtoUSART_Buffer, ReceivedByte);
    }

    uint16_t BufferCount = RingBuffer_GetCount(&USARTtoUSB_Buffer);
    if ((TIFR0 & (1 << TOV0)) || (BufferCount > (uint8_t)(sizeof(USARTtoUSB_Buffer_Data) * .75))) {
      TIFR0 |= (1 << TOV0);
      while (BufferCount--) {
        if (CDC_Device_SendByte(&VirtualSerial1_CDC_Interface, RingBuffer_Peek(&USARTtoUSB_Buffer)) != ENDPOINT_READYWAIT_NoError)
          break;

        RingBuffer_Remove(&USARTtoUSB_Buffer);
      }
    }

    if (!(RingBuffer_IsEmpty(&USBtoUSART_Buffer)))
      Serial_SendByte(RingBuffer_Remove(&USBtoUSART_Buffer));
}
/* end LUFA CDC-ACM specific definitions*/




void stdio_init(void)
{
  CDC_Device_CreateStream(&VirtualSerial1_CDC_Interface, &usb1_stream_);
  CDC_Device_CreateStream(&VirtualSerial2_CDC_Interface, &usb2_stream_);
  stdin = &usb2_stream_;
  stdout = &usb2_stream_;
  stderr = &usb2_stream_;
}
