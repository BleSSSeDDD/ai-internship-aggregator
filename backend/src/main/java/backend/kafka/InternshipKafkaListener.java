package backend.kafka;

import backend.service.InternshipService;
import com.aggregator.internship.CompanyInternship;
import com.google.protobuf.InvalidProtocolBufferException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.List;

@Component
@RequiredArgsConstructor
@Slf4j
public class InternshipKafkaListener {

    private final InternshipService internshipService;


    @KafkaListener(
            topics = "internships",
            groupId = "internship-db-consumer",
            containerFactory = "kafkaBatchListenerContainerFactory"
    )
    public void listenBatch(List<byte[]> payloads) {
        List<CompanyInternship> internships = new ArrayList<>();

        for (byte[] payload : payloads) {
            try {
                CompanyInternship internship = CompanyInternship.parseFrom(payload);

                if (internship == null || internship.getPositionName().isEmpty()) {
                    log.warn("Received invalid internship data: {}", internship);
                    continue;
                }

                internships.add(internship);

            } catch (InvalidProtocolBufferException e) {
                log.error("Failed to parse Protobuf message", e);
            } catch (Exception e) {
                log.error("Unexpected error processing message", e);
            }
        }

        if (!internships.isEmpty()) {
            internshipService.saveBatch(internships);
            log.debug("Successfully processed batch of {} internships", internships.size());
        }
    }
}