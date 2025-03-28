CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE public.users (
                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                              full_name VARCHAR(100) ,
                              phone_number VARCHAR(20) ,
                              username VARCHAR(50) UNIQUE NOT NULL,
                              email VARCHAR(100) UNIQUE,
                              password_hash VARCHAR(255) NOT NULL,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.doctors (
                                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                first_name text NOT NULL,
                                last_name text NOT NULL,
                                slug text NOT NULL,
                                image character varying(255) NOT NULL,
                                experience integer NOT NULL,
                                specialization text NOT NULL,
                                treatment_profile text NOT NULL,
                                professional_activity text NOT NULL,
                                working_hours jsonb NOT NULL,
                                created_at timestamp(0) without time zone,
                                updated_at timestamp(0) without time zone
);

CREATE TYPE appointment_status AS ENUM ('pending', 'confirmed', 'canceled');
CREATE TABLE public.appointments (
                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                              user_id UUID REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE ,
                              doctor_id UUID REFERENCES doctors(id) ON DELETE CASCADE ON UPDATE CASCADE ,
                              appointment_time TIMESTAMP NOT NULL,
                              status appointment_status DEFAULT 'pending',
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              CONSTRAINT unique_appointment UNIQUE (doctor_id, appointment_time)
);

CREATE TABLE public.notifications (
                                      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                      appointment_id UUID REFERENCES public.appointments(id) ON DELETE CASCADE,
                                      user_id UUID REFERENCES public.users(id) ON DELETE CASCADE,
                                      doctor_id UUID REFERENCES public.doctors(id) ON DELETE CASCADE,
                                      scheduled_at TIMESTAMP NOT NULL,
                                      sent_at TIMESTAMP,
                                      message TEXT NOT NULL,
                                      status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed'))
);