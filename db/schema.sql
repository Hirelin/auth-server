-- CREATE DATABASE "hirelin-db";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CreateEnum
CREATE TYPE "AccountType" AS ENUM ('oauth', 'email', 'credentials');

-- CreateEnum
CREATE TYPE "ApplicationStatus" AS ENUM ('pending', 'rejected', 'accepted', 'training', 'hired');

-- CreateEnum
CREATE TYPE "UploadType" AS ENUM ('resume', 'jobDescription', 'requirements');

-- CreateEnum
CREATE TYPE "JobType" AS ENUM ('fullTime', 'partTime', 'contract', 'internship', 'freelance', 'temporary', 'volunteer', 'remote', 'onSite', 'hybrid');

-- CreateEnum
CREATE TYPE "JobStatus" AS ENUM ('upcoming', 'open', 'closed', 'cancelled', 'completed');

-- CreateTable
CREATE TABLE "User" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "name" TEXT,
    "email" TEXT NOT NULL,
    "image" TEXT NOT NULL DEFAULT '/images/profile-default.jpg',
    "emailVerified" TIMESTAMP(3),
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "User_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Account" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "user_id" UUID NOT NULL,
    "type" "AccountType" NOT NULL DEFAULT 'oauth',
    "provider" TEXT NOT NULL,
    "provider_account_id" TEXT NOT NULL,
    "access_token" TEXT,
    "refresh_token" TEXT,
    "expires_at" INTEGER,
    "token_type" TEXT,
    "id_token" TEXT,
    "session_state" TEXT,
    "scope" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Account_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Session" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "session_token" TEXT NOT NULL,
    "refresh_token" TEXT NOT NULL,
    "expires_at" TIMESTAMP(3) NOT NULL,
    "user_id" UUID NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Session_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "VerificationToken" (
    "identifier" TEXT NOT NULL,
    "token" TEXT NOT NULL,
    "expires_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "VerificationToken_pkey" PRIMARY KEY ("identifier","token")
);

-- CreateTable
CREATE TABLE "Recruiter" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "name" TEXT NOT NULL,
    "position" TEXT NOT NULL,
    "organization" TEXT NOT NULL,
    "phone" TEXT,
    "address" TEXT,
    "user_id" UUID NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Recruiter_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "JobOpening" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "title" TEXT NOT NULL,
    "company" TEXT NOT NULL,
    "location" TEXT,
    "type" "JobType" NOT NULL,
    "description" TEXT NOT NULL,
    "contact" TEXT NOT NULL,
    "address" TEXT,
    "status" "JobStatus" NOT NULL DEFAULT 'upcoming',
    "deadline" TIMESTAMP(3),
    "start_date" TIMESTAMP(3),
    "end_date" TIMESTAMP(3),
    "requirements_file_id" UUID,
    "parsed_requirements" TEXT,
    "training_id" UUID,
    "recruiter_id" UUID NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "JobOpening_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Application" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "status" "ApplicationStatus" NOT NULL DEFAULT 'pending',
    "job_opening_id" UUID NOT NULL,
    "user_id" UUID NOT NULL,
    "parsed_resume" TEXT,
    "layout_score" DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    "content_score" DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    "skill_gap" TEXT,
    "resume_id" UUID,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Application_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Uploads" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "uploadType" "UploadType" NOT NULL,
    "name" TEXT NOT NULL,
    "file_type" TEXT NOT NULL,
    "url" TEXT NOT NULL,
    "user_id" UUID NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Uploads_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Training" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "topics" TEXT NOT NULL,
    "start_date" TIMESTAMP(3),
    "end_date" TIMESTAMP(3),
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Training_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "LearningPlan" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "plan_details" JSONB NOT NULL,
    "application_id" UUID NOT NULL,
    "training_id" UUID NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "LearningPlan_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Assessment" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "title" TEXT NOT NULL,
    "description" TEXT,
    "questions" JSONB NOT NULL,
    "learning_plan_id" UUID NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Assessment_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "User_email_key" ON "User"("email");

-- CreateIndex
CREATE UNIQUE INDEX "Session_session_token_key" ON "Session"("session_token");

-- CreateIndex
CREATE UNIQUE INDEX "Session_refresh_token_key" ON "Session"("refresh_token");

-- CreateIndex
CREATE INDEX "VerificationToken_identifier_token_idx" ON "VerificationToken"("identifier", "token");

-- CreateIndex
CREATE UNIQUE INDEX "Recruiter_user_id_key" ON "Recruiter"("user_id");

-- CreateIndex
CREATE INDEX "JobOpening_recruiter_id_idx" ON "JobOpening"("recruiter_id");

-- CreateIndex
CREATE INDEX "Application_job_opening_id_idx" ON "Application"("job_opening_id");

-- CreateIndex
CREATE INDEX "Application_user_id_idx" ON "Application"("user_id");

-- CreateIndex
CREATE UNIQUE INDEX "LearningPlan_application_id_key" ON "LearningPlan"("application_id");

-- CreateIndex
CREATE INDEX "Assessment_learning_plan_id_idx" ON "Assessment"("learning_plan_id");

-- AddForeignKey
ALTER TABLE "Account" ADD CONSTRAINT "Account_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Session" ADD CONSTRAINT "Session_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Recruiter" ADD CONSTRAINT "Recruiter_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "JobOpening" ADD CONSTRAINT "JobOpening_requirements_file_id_fkey" FOREIGN KEY ("requirements_file_id") REFERENCES "Uploads"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "JobOpening" ADD CONSTRAINT "JobOpening_training_id_fkey" FOREIGN KEY ("training_id") REFERENCES "Training"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "JobOpening" ADD CONSTRAINT "JobOpening_recruiter_id_fkey" FOREIGN KEY ("recruiter_id") REFERENCES "Recruiter"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Application" ADD CONSTRAINT "Application_job_opening_id_fkey" FOREIGN KEY ("job_opening_id") REFERENCES "JobOpening"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Application" ADD CONSTRAINT "Application_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Application" ADD CONSTRAINT "Application_resume_id_fkey" FOREIGN KEY ("resume_id") REFERENCES "Uploads"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Uploads" ADD CONSTRAINT "Uploads_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "LearningPlan" ADD CONSTRAINT "LearningPlan_application_id_fkey" FOREIGN KEY ("application_id") REFERENCES "Application"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "LearningPlan" ADD CONSTRAINT "LearningPlan_training_id_fkey" FOREIGN KEY ("training_id") REFERENCES "Training"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Assessment" ADD CONSTRAINT "Assessment_learning_plan_id_fkey" FOREIGN KEY ("learning_plan_id") REFERENCES "LearningPlan"("id") ON DELETE CASCADE ON UPDATE CASCADE;


CREATE OR REPLACE FUNCTION set_current_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER updated_at_User_trigger
BEFORE UPDATE ON "User"
FOR EACH ROW
EXECUTE FUNCTION set_current_timestamp();

CREATE TRIGGER updated_at_Session_trigger
BEFORE UPDATE ON "Session"
FOR EACH ROW
EXECUTE FUNCTION set_current_timestamp();

CREATE TRIGGER updated_at_Recruiter_trigger
BEFORE UPDATE ON "Recruiter"
FOR EACH ROW
EXECUTE FUNCTION set_current_timestamp();

CREATE TRIGGER updated_at_Uploads_trigger
BEFORE UPDATE ON "Uploads"
FOR EACH ROW
EXECUTE FUNCTION set_current_timestamp();

CREATE TRIGGER updated_at_Training_trigger
BEFORE UPDATE ON "Training"
FOR EACH ROW
EXECUTE FUNCTION set_current_timestamp();

CREATE TRIGGER updated_at_LearningPlan_trigger
BEFORE UPDATE ON "LearningPlan"
FOR EACH ROW
EXECUTE FUNCTION set_current_timestamp();

CREATE TRIGGER updated_at_Assessment_trigger
BEFORE UPDATE ON "Assessment"
FOR EACH ROW
EXECUTE FUNCTION set_current_timestamp();

